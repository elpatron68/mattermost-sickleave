// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import manifest from 'manifest';

import type {GlobalState} from '@mattermost/types/store';

import type {SickLeaveContext, SickLeaveVariant} from 'types';

type PluginState = {
    modalVisible: boolean;
    menuVisible: boolean;
    menuChannelId: string;
    menuTeamId: string;
    menuLoading: boolean;
    menuEnding: boolean;
    menuError: string;
    variant: SickLeaveVariant | '';
    channelId: string;
    teamId: string;
    context: SickLeaveContext | null;
    submitting: boolean;
    fieldErrors: Record<string, string>;
    generalError: string;
};

function getPluginState(state: GlobalState): PluginState | undefined {
    return (state as GlobalState & Record<string, PluginState>)['plugins-' + manifest.id];
}

export function isSickLeaveModalVisible(state: GlobalState): boolean {
    return getPluginState(state)?.modalVisible ?? false;
}

export function isSickLeaveMenuVisible(state: GlobalState): boolean {
    return getPluginState(state)?.menuVisible ?? false;
}

export function sickLeaveMenuChannelId(state: GlobalState): string {
    return getPluginState(state)?.menuChannelId ?? '';
}

export function sickLeaveMenuTeamId(state: GlobalState): string {
    return getPluginState(state)?.menuTeamId ?? '';
}

export function sickLeaveMenuLoading(state: GlobalState): boolean {
    return getPluginState(state)?.menuLoading ?? false;
}

export function sickLeaveMenuEnding(state: GlobalState): boolean {
    return getPluginState(state)?.menuEnding ?? false;
}

export function sickLeaveMenuError(state: GlobalState): string {
    return getPluginState(state)?.menuError ?? '';
}

export function sickLeaveVariant(state: GlobalState): SickLeaveVariant | '' {
    return getPluginState(state)?.variant ?? '';
}

export function sickLeaveChannelId(state: GlobalState): string {
    return getPluginState(state)?.channelId ?? '';
}

export function sickLeaveTeamId(state: GlobalState): string {
    return getPluginState(state)?.teamId ?? '';
}

export function sickLeaveContext(state: GlobalState): SickLeaveContext | null {
    return getPluginState(state)?.context ?? null;
}

export function sickLeaveSubmitting(state: GlobalState): boolean {
    return getPluginState(state)?.submitting ?? false;
}

export function sickLeaveFieldErrors(state: GlobalState): Record<string, string> {
    return getPluginState(state)?.fieldErrors ?? {};
}

export function sickLeaveGeneralError(state: GlobalState): string {
    return getPluginState(state)?.generalError ?? '';
}

export function sickLeaveCommandTrigger(state: GlobalState): string {
    return getPluginState(state)?.context?.command_trigger || 'sick-leave';
}
