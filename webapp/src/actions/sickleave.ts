// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {Client4} from 'mattermost-redux/client';
import {getCurrentUserLocale} from 'mattermost-redux/selectors/entities/i18n';

import type {GlobalState} from '@mattermost/types/store';

import {
    CLOSE_SICK_LEAVE_MENU,
    CLOSE_SICK_LEAVE_MODAL,
    OPEN_SICK_LEAVE_MENU,
    OPEN_SICK_LEAVE_MODAL,
    SET_SICK_LEAVE_CONTEXT,
    SET_SICK_LEAVE_ERRORS,
    SET_SICK_LEAVE_MENU_ENDING,
    SET_SICK_LEAVE_MENU_ERROR,
    SET_SICK_LEAVE_MENU_LOADING,
    SET_SICK_LEAVE_SUBMITTING,
} from 'action_types';
import {endSickLeave, fetchSickLeaveContext, submitSickLeaveDialog} from 'client';
import {
    sickLeaveChannelId,
    sickLeaveMenuChannelId,
    sickLeaveMenuTeamId,
    sickLeaveTeamId,
    sickLeaveVariant,
} from 'selectors';
import type {MenuAction, SickLeaveVariant} from 'types';

import de from '../../i18n/de.json';
import en from '../../i18n/en.json';

type Dispatch = (action: {type: string; [key: string]: unknown}) => void;

function getTranslationsForLocale(locale: string): Record<string, string> {
    switch (locale) {
    case 'de':
        return de;
    default:
        return en;
    }
}

export function confirmEndCase(getState: () => GlobalState): boolean {
    const locale = getCurrentUserLocale(getState());
    const translations = getTranslationsForLocale(locale);
    return window.confirm(translations['menu.end.confirm'] || 'Close your active sick leave case? HR will be notified.');
}

async function showEphemeralMessage(channelId: string, message: string, getState: () => GlobalState): Promise<void> {
    const userId = getState().entities.users.currentUserId;
    await Client4.createPost({
        user_id: userId,
        channel_id: channelId,
        message,
        type: 'system_ephemeral',
    });
}

export const closeSickLeaveModal = () => (dispatch: Dispatch) => {
    dispatch({type: CLOSE_SICK_LEAVE_MODAL});
};

export const closeSickLeaveMenu = () => (dispatch: Dispatch) => {
    dispatch({type: CLOSE_SICK_LEAVE_MENU});
};

export const openSickLeaveMenu = (channelId: string, teamId: string) => async (dispatch: Dispatch, getState: () => GlobalState) => {
    dispatch({type: OPEN_SICK_LEAVE_MENU, channelId, teamId});
    dispatch({type: SET_SICK_LEAVE_MENU_ERROR, error: ''});

    try {
        const context = await fetchSickLeaveContext();
        dispatch({type: SET_SICK_LEAVE_CONTEXT, context});
        dispatch({type: SET_SICK_LEAVE_MENU_LOADING, loading: false});
    } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to open sick leave menu';
        dispatch({type: SET_SICK_LEAVE_MENU_LOADING, loading: false});
        dispatch({type: CLOSE_SICK_LEAVE_MENU});
        await showEphemeralMessage(channelId, message, getState);
    }
};

export const openSickLeaveModal = (variant: SickLeaveVariant, channelId: string, teamId: string) => async (dispatch: Dispatch, getState: () => GlobalState) => {
    try {
        const context = await fetchSickLeaveContext();
        dispatch({type: SET_SICK_LEAVE_CONTEXT, context});
        dispatch({type: OPEN_SICK_LEAVE_MODAL, variant, channelId, teamId});
    } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to open sick leave dialog';
        await showEphemeralMessage(channelId, message, getState);
    }
};

export const endSickLeaveCase = (channelId: string) => async (dispatch: Dispatch, getState: () => GlobalState) => {
    dispatch({type: SET_SICK_LEAVE_MENU_ENDING, ending: true});
    dispatch({type: SET_SICK_LEAVE_MENU_ERROR, error: ''});

    try {
        const response = await endSickLeave(channelId);
        dispatch({type: CLOSE_SICK_LEAVE_MENU});
        await showEphemeralMessage(channelId, response.message, getState);
    } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to close sick leave case';
        dispatch({type: SET_SICK_LEAVE_MENU_ERROR, error: message});
        dispatch({type: SET_SICK_LEAVE_MENU_ENDING, ending: false});
    }
};

export const selectSickLeaveMenuAction = (action: MenuAction) => async (dispatch: Dispatch, getState: () => GlobalState) => {
    const state = getState();
    const channelId = sickLeaveMenuChannelId(state);
    let teamId = sickLeaveMenuTeamId(state);
    if (!teamId && channelId) {
        teamId = state.entities.channels.channels[channelId]?.team_id || '';
    }

    if (action === 'end') {
        if (!confirmEndCase(getState)) {
            return;
        }
        await endSickLeaveCase(channelId)(dispatch, getState);
        return;
    }

    if (action === 'status') {
        return;
    }

    dispatch({type: CLOSE_SICK_LEAVE_MENU});
    await openSickLeaveModal(action, channelId, teamId)(dispatch, getState);
};

type SubmitPayload = {
    startDate?: string;
    expectedEndDate?: string;
    auCertificate?: string;
};

export const submitSickLeaveModal = (payload: SubmitPayload) => async (dispatch: Dispatch, getState: () => GlobalState) => {
    const state = getState();
    const variant = sickLeaveVariant(state);
    const channelId = sickLeaveChannelId(state);
    if (!variant || !channelId) {
        return;
    }

    let teamId = sickLeaveTeamId(state);
    if (!teamId) {
        teamId = state.entities.channels.channels[channelId]?.team_id || '';
    }

    const submission: Record<string, string> = {};
    if (payload.startDate) {
        submission.start_date = payload.startDate;
    }
    if (payload.expectedEndDate) {
        submission.expected_end_date = payload.expectedEndDate;
    }
    if (payload.auCertificate) {
        submission.au_certificate = payload.auCertificate;
    }

    dispatch({type: SET_SICK_LEAVE_SUBMITTING, submitting: true});
    dispatch({type: SET_SICK_LEAVE_ERRORS, fieldErrors: {}, generalError: ''});

    try {
        const response = await submitSickLeaveDialog(variant, channelId, teamId, submission);
        const errors = response.errors || {};
        const fieldErrors = {...errors};
        const generalError = fieldErrors[''] || '';
        delete fieldErrors[''];

        if (generalError || Object.keys(fieldErrors).length > 0) {
            dispatch({type: SET_SICK_LEAVE_ERRORS, fieldErrors, generalError});
            dispatch({type: SET_SICK_LEAVE_SUBMITTING, submitting: false});
            return;
        }

        dispatch({type: CLOSE_SICK_LEAVE_MODAL});
    } catch (error) {
        const message = error instanceof Error ? error.message : 'Submission failed';
        dispatch({type: SET_SICK_LEAVE_ERRORS, fieldErrors: {}, generalError: message});
        dispatch({type: SET_SICK_LEAVE_SUBMITTING, submitting: false});
    }
};

export function parseSickLeaveCommand(message: string): string | null {
    const trimmed = message.trim();
    const match = /^\/sick-leave(?:\s+(\S+))?$/i.exec(trimmed);
    if (!match) {
        return null;
    }
    return (match[1] || 'help').toLowerCase();
}
