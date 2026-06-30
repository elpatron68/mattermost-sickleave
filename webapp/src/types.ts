// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export type SickLeaveVariant = 'start' | 'update' | 'extend';

export type MenuAction = 'start' | 'update' | 'extend' | 'end' | 'status';

export type ActiveRecord = {
    start_date: string;
    expected_end_date?: string;
    status: string;
};

export type SickLeaveContext = {
    active?: ActiveRecord;
    max_backdate_days: number;
    command_trigger: string;
};
