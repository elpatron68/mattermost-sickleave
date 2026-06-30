// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import ConnectedSickLeaveMenu from 'components/sick_leave_menu';
import ConnectedSickLeaveModal from 'components/sick_leave_modal';

export default function Root(): React.ReactElement {
    return (
        <>
            <ConnectedSickLeaveMenu/>
            <ConnectedSickLeaveModal/>
        </>
    );
}
