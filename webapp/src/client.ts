// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {Client4} from 'mattermost-redux/client';

import type {SickLeaveContext, SickLeaveVariant} from 'types';
import {getPluginURL} from 'utils';

type SubmitDialogResponse = {
    errors?: Record<string, string>;
};

const callbackIDs: Record<SickLeaveVariant, string> = {
    start: 'sickleave_start',
    update: 'sickleave_update',
    extend: 'sickleave_extend',
};

export async function fetchSickLeaveContext(): Promise<SickLeaveContext> {
    const response = await fetch(`${getPluginURL()}/api/v1/context`, {
        headers: Client4.getOptions({method: 'get'}).headers,
    });

    if (!response.ok) {
        throw new Error('Failed to load sick leave context');
    }

    return response.json() as Promise<SickLeaveContext>;
}

export async function submitSickLeaveDialog(
    variant: SickLeaveVariant,
    channelId: string,
    teamId: string,
    submission: Record<string, string>,
): Promise<SubmitDialogResponse> {
    const response = await fetch(`${getPluginURL()}/api/v1/dialog/submit`, {
        method: 'POST',
        headers: {
            ...Client4.getOptions({method: 'post'}).headers,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            callback_id: callbackIDs[variant],
            channel_id: channelId,
            team_id: teamId,
            submission,
        }),
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Submission failed');
    }

    return response.json() as Promise<SubmitDialogResponse>;
}

type EndResponse = {
    message: string;
};

export async function endSickLeave(channelId: string): Promise<EndResponse> {
    const response = await fetch(`${getPluginURL()}/api/v1/end`, {
        method: 'POST',
        headers: {
            ...Client4.getOptions({method: 'post'}).headers,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({channel_id: channelId}),
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Failed to close sick leave case');
    }

    return response.json() as Promise<EndResponse>;
}
