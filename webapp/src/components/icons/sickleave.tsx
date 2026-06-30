// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

const CROSS_COLOR = '#D32F2F';

export default function SickLeaveIcon(): React.ReactElement {
    return (
        <svg
            width='18'
            height='18'
            viewBox='0 0 24 24'
            aria-hidden='true'
            role='img'
        >
            <rect x='10' y='3' width='4' height='18' rx='1' fill={CROSS_COLOR}/>
            <rect x='3' y='10' width='18' height='4' rx='1' fill={CROSS_COLOR}/>
        </svg>
    );
}
