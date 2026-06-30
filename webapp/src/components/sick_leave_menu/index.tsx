// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators, type Dispatch} from 'redux';

import type {GlobalState} from '@mattermost/types/store';

import {getTheme} from 'mattermost-redux/selectors/entities/preferences';

import {closeSickLeaveMenu, selectSickLeaveMenuAction} from 'actions/sickleave';
import {
    isSickLeaveMenuVisible,
    sickLeaveContext,
    sickLeaveMenuEnding,
    sickLeaveMenuError,
    sickLeaveMenuLoading,
} from 'selectors';

import SickLeaveMenu from './sick_leave_menu';

function mapStateToProps(state: GlobalState) {
    return {
        visible: isSickLeaveMenuVisible(state),
        context: sickLeaveContext(state),
        loading: sickLeaveMenuLoading(state),
        ending: sickLeaveMenuEnding(state),
        error: sickLeaveMenuError(state),
        theme: getTheme(state),
    };
}

function mapDispatchToProps(dispatch: Dispatch) {
    return bindActionCreators({
        onClose: closeSickLeaveMenu,
        onSelectAction: selectSickLeaveMenuAction,
    }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(SickLeaveMenu);
