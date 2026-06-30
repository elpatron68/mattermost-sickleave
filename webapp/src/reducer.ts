// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

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
import {combineReducers} from 'redux';

import type {SickLeaveContext, SickLeaveVariant} from 'types';

type OpenModalAction = {
    type: typeof OPEN_SICK_LEAVE_MODAL;
    variant: SickLeaveVariant;
    channelId: string;
    teamId: string;
};

type SetContextAction = {
    type: typeof SET_SICK_LEAVE_CONTEXT;
    context: SickLeaveContext;
};

type SetErrorsAction = {
    type: typeof SET_SICK_LEAVE_ERRORS;
    fieldErrors: Record<string, string>;
    generalError: string;
};

type OpenMenuAction = {
    type: typeof OPEN_SICK_LEAVE_MENU;
    channelId: string;
    teamId: string;
};

type SetMenuErrorAction = {
    type: typeof SET_SICK_LEAVE_MENU_ERROR;
    error: string;
};

const menuVisible = (state = false, action: {type: string}) => {
    switch (action.type) {
    case OPEN_SICK_LEAVE_MENU:
        return true;
    case CLOSE_SICK_LEAVE_MENU:
        return false;
    default:
        return state;
    }
};

const menuChannelId = (state = '', action: {type: string}) => {
    switch (action.type) {
    case OPEN_SICK_LEAVE_MENU:
        return (action as OpenMenuAction).channelId;
    case CLOSE_SICK_LEAVE_MENU:
        return '';
    default:
        return state;
    }
};

const menuTeamId = (state = '', action: {type: string}) => {
    switch (action.type) {
    case OPEN_SICK_LEAVE_MENU:
        return (action as OpenMenuAction).teamId;
    case CLOSE_SICK_LEAVE_MENU:
        return '';
    default:
        return state;
    }
};

const menuLoading = (state = false, action: {type: string; loading?: boolean}) => {
    switch (action.type) {
    case SET_SICK_LEAVE_MENU_LOADING:
        return action.loading ?? false;
    case OPEN_SICK_LEAVE_MENU:
        return true;
    case CLOSE_SICK_LEAVE_MENU:
        return false;
    default:
        return state;
    }
};

const menuEnding = (state = false, action: {type: string; ending?: boolean}) => {
    switch (action.type) {
    case SET_SICK_LEAVE_MENU_ENDING:
        return action.ending ?? false;
    case CLOSE_SICK_LEAVE_MENU:
        return false;
    default:
        return state;
    }
};

const menuError = (state = '', action: {type: string}) => {
    switch (action.type) {
    case SET_SICK_LEAVE_MENU_ERROR:
        return (action as SetMenuErrorAction).error;
    case OPEN_SICK_LEAVE_MENU:
    case CLOSE_SICK_LEAVE_MENU:
        return '';
    default:
        return state;
    }
};

const modalVisible = (state = false, action: {type: string}) => {
    switch (action.type) {
    case OPEN_SICK_LEAVE_MODAL:
        return true;
    case CLOSE_SICK_LEAVE_MODAL:
        return false;
    default:
        return state;
    }
};

const variant = (state: SickLeaveVariant | '' = '', action: {type: string}) => {
    switch (action.type) {
    case OPEN_SICK_LEAVE_MODAL:
        return (action as OpenModalAction).variant;
    case CLOSE_SICK_LEAVE_MODAL:
        return '';
    default:
        return state;
    }
};

const channelId = (state = '', action: {type: string}) => {
    switch (action.type) {
    case OPEN_SICK_LEAVE_MODAL:
        return (action as OpenModalAction).channelId;
    case CLOSE_SICK_LEAVE_MODAL:
        return '';
    default:
        return state;
    }
};

const teamId = (state = '', action: {type: string}) => {
    switch (action.type) {
    case OPEN_SICK_LEAVE_MODAL:
        return (action as OpenModalAction).teamId;
    case CLOSE_SICK_LEAVE_MODAL:
        return '';
    default:
        return state;
    }
};

const context = (state: SickLeaveContext | null = null, action: {type: string}) => {
    switch (action.type) {
    case SET_SICK_LEAVE_CONTEXT:
        return (action as SetContextAction).context;
    case CLOSE_SICK_LEAVE_MODAL:
        return null;
    case CLOSE_SICK_LEAVE_MENU:
        return state;
    default:
        return state;
    }
};

const submitting = (state = false, action: {type: string; submitting?: boolean}) => {
    switch (action.type) {
    case SET_SICK_LEAVE_SUBMITTING:
        return action.submitting ?? false;
    case CLOSE_SICK_LEAVE_MODAL:
        return false;
    default:
        return state;
    }
};

const fieldErrors = (state: Record<string, string> = {}, action: {type: string}) => {
    switch (action.type) {
    case SET_SICK_LEAVE_ERRORS:
        return (action as SetErrorsAction).fieldErrors;
    case CLOSE_SICK_LEAVE_MODAL:
    case OPEN_SICK_LEAVE_MODAL:
        return {};
    default:
        return state;
    }
};

const generalError = (state = '', action: {type: string}) => {
    switch (action.type) {
    case SET_SICK_LEAVE_ERRORS:
        return (action as SetErrorsAction).generalError;
    case CLOSE_SICK_LEAVE_MODAL:
    case OPEN_SICK_LEAVE_MODAL:
        return '';
    default:
        return state;
    }
};

export default combineReducers({
    modalVisible,
    menuVisible,
    menuChannelId,
    menuTeamId,
    menuLoading,
    menuEnding,
    menuError,
    variant,
    channelId,
    teamId,
    context,
    submitting,
    fieldErrors,
    generalError,
});
