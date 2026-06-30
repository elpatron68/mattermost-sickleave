// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {MenuAction, SickLeaveContext} from 'types';

export function getAvailableMenuActions(context: SickLeaveContext | null): MenuAction[] {
    if (!context?.active) {
        return ['start', 'status'];
    }

    switch (context.active.status) {
    case 'reported':
        return ['update', 'end', 'status'];
    case 'updated':
    case 'extended':
        return ['extend', 'end', 'status'];
    default:
        return ['status'];
    }
}
