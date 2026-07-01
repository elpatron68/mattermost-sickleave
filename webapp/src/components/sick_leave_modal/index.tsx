// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators, type Dispatch} from 'redux';

import type {GlobalState} from '@mattermost/types/store';

import {getTheme} from 'mattermost-redux/selectors/entities/preferences';
import {getCurrentUserLocale} from 'mattermost-redux/selectors/entities/i18n';

import {closeSickLeaveModal, submitSickLeaveModal} from 'actions/sickleave';
import {
    isSickLeaveModalVisible,
    sickLeaveContext,
    sickLeaveFieldErrors,
    sickLeaveGeneralError,
    sickLeaveSubmitting,
    sickLeaveVariant,
} from 'selectors';

import SickLeaveModal from './sick_leave_modal';

function mapStateToProps(state: GlobalState) {
    return {
        visible: isSickLeaveModalVisible(state),
        variant: sickLeaveVariant(state),
        locale: getCurrentUserLocale(state),
        context: sickLeaveContext(state),
        submitting: sickLeaveSubmitting(state),
        fieldErrors: sickLeaveFieldErrors(state),
        generalError: sickLeaveGeneralError(state),
        theme: getTheme(state),
    };
}

function mapDispatchToProps(dispatch: Dispatch) {
    return bindActionCreators({
        onClose: closeSickLeaveModal,
        onSubmit: submitSickLeaveModal,
    }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(SickLeaveModal);
