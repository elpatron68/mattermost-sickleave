// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage} from 'react-intl';
import {getAvailableMenuActions} from 'utils/menu';

import {changeOpacity} from 'mattermost-redux/utils/theme_utils';

import type {MenuAction, SickLeaveContext} from 'types';

import './sick_leave_menu.css';

type Props = {
    visible: boolean;
    context: SickLeaveContext | null;
    loading: boolean;
    ending: boolean;
    error: string;
    onClose: () => void;
    onSelectAction: (action: MenuAction) => void;
    theme: {
        centerChannelBg: string;
        centerChannelColor: string;
        errorTextColor: string;
    };
};

export default class SickLeaveMenu extends React.PureComponent<Props> {
    private labelForAction(action: MenuAction): React.ReactNode {
        switch (action) {
        case 'start':
            return (
                <FormattedMessage
                    id='menu.action.start'
                    defaultMessage='Report sick leave'
                />
            );
        case 'update':
            return (
                <FormattedMessage
                    id='menu.action.update'
                    defaultMessage='Update sick leave'
                />
            );
        case 'extend':
            return (
                <FormattedMessage
                    id='menu.action.extend'
                    defaultMessage='Extend sick leave'
                />
            );
        case 'end':
            return (
                <FormattedMessage
                    id='menu.action.end'
                    defaultMessage='Close sick leave case'
                />
            );
        case 'status':
            return (
                <FormattedMessage
                    id='menu.action.status'
                    defaultMessage='Show status'
                />
            );
        default:
            return null;
        }
    }

    private renderStatus(): React.ReactNode {
        const {context} = this.props;
        if (!context?.active) {
            return (
                <p className='sickleave-menu__status'>
                    <FormattedMessage
                        id='menu.status.none'
                        defaultMessage='You do not have an active sick leave report.'
                    />
                </p>
            );
        }

        const active = context.active;
        if (!active.expected_end_date) {
            return (
                <p className='sickleave-menu__status'>
                    <FormattedMessage
                        id='menu.status.active'
                        defaultMessage='Active sick leave since {startDate}. Status: {status}.'
                        values={{startDate: active.start_date, status: active.status}}
                    />
                </p>
            );
        }

        return (
            <p className='sickleave-menu__status'>
                <FormattedMessage
                    id='menu.status.active_with_end'
                    defaultMessage='Active sick leave since {startDate}, expected return {endDate}. Status: {status}.'
                    values={{
                        startDate: active.start_date,
                        endDate: active.expected_end_date,
                        status: active.status,
                    }}
                />
            </p>
        );
    }

    render(): React.ReactNode {
        if (!this.props.visible) {
            return null;
        }

        const actions = getAvailableMenuActions(this.props.context);
        const style = getStyle(this.props.theme);
        const busy = this.props.loading || this.props.ending;

        return (
            <div
                className='sickleave-modal__overlay'
                style={style.overlay}
                onClick={this.props.onClose}
            >
                <div
                    className='sickleave-menu__panel'
                    style={style.panel}
                    onClick={(event) => event.stopPropagation()}
                    role='dialog'
                    aria-modal='true'
                >
                    <h2 className='sickleave-menu__title'>
                        <FormattedMessage
                            id='menu.title'
                            defaultMessage='Sick leave'
                        />
                    </h2>
                    {this.renderStatus()}
                    <div className='sickleave-menu__actions'>
                        {actions.filter((action) => action !== 'status').map((action) => (
                            <button
                                key={action}
                                type='button'
                                className={`sickleave-menu__action${action === 'end' ? ' sickleave-menu__action--danger' : ''}`}
                                onClick={() => this.props.onSelectAction(action)}
                                disabled={busy}
                            >
                                {this.labelForAction(action)}
                            </button>
                        ))}
                    </div>
                    {this.props.error && (
                        <div className='sickleave-menu__error'>{this.props.error}</div>
                    )}
                    <div className='sickleave-menu__footer'>
                        <button
                            type='button'
                            className='sickleave-menu__button'
                            onClick={this.props.onClose}
                            disabled={busy}
                        >
                            <FormattedMessage
                                id='menu.close'
                                defaultMessage='Close'
                            />
                        </button>
                    </div>
                </div>
            </div>
        );
    }
}

function getStyle(theme: Props['theme']) {
    return {
        overlay: {
            position: 'fixed' as const,
            display: 'flex',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            zIndex: 2000,
            alignItems: 'center',
            justifyContent: 'center',
            backgroundColor: 'rgba(0, 0, 0, 0.5)',
        },
        panel: {
            backgroundColor: theme.centerChannelBg,
            color: theme.centerChannelColor,
            border: `1px solid ${changeOpacity(theme.centerChannelColor, 0.1)}`,
        },
    };
}
