// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {FormattedMessage} from 'react-intl';
import {changeOpacity} from 'mattermost-redux/utils/theme_utils';

import type {SickLeaveContext, SickLeaveVariant} from 'types';
import {addDays, formatISODate, parseISODate} from 'utils';

import './sick_leave_modal.css';

type Props = {
    visible: boolean;
    variant: SickLeaveVariant | '';
    context: SickLeaveContext | null;
    submitting: boolean;
    fieldErrors: Record<string, string>;
    generalError: string;
    onClose: () => void;
    onSubmit: (payload: {startDate?: string; expectedEndDate?: string; auCertificate?: string}) => void;
    theme: {
        centerChannelBg: string;
        centerChannelColor: string;
        linkColor: string;
        errorTextColor: string;
    };
};

type State = {
    startDate: string;
    expectedEndDate: string;
    auCertificate: string;
};

export default class SickLeaveModal extends React.PureComponent<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = this.initialState(props);
    }

    componentDidUpdate(prevProps: Props): void {
        if (prevProps.visible !== this.props.visible && this.props.visible) {
            this.setState(this.initialState(this.props));
        }
    }

    private initialState(props: Props): State {
        const today = formatISODate(new Date());
        return {
            startDate: today,
            expectedEndDate: props.context?.active?.expected_end_date || today,
            auCertificate: props.variant === 'extend' ? 'unchanged' : '',
        };
    }

    private getDateBounds(): {minDate?: string; maxDate?: string} {
        const {variant, context} = this.props;
        const today = new Date();

        if (variant === 'start') {
            const maxBackdate = context?.max_backdate_days || 3;
            return {
                minDate: formatISODate(addDays(today, -maxBackdate)),
                maxDate: formatISODate(today),
            };
        }

        if (variant === 'update' && context?.active?.start_date) {
            return {minDate: context.active.start_date};
        }

        if (variant === 'extend' && context?.active?.expected_end_date) {
            const current = parseISODate(context.active.expected_end_date);
            if (current) {
                return {minDate: formatISODate(addDays(current, 1))};
            }
        }

        return {};
    }

    private handleSubmit = (event: React.FormEvent): void => {
        event.preventDefault();
        const {variant} = this.props;
        const {startDate, expectedEndDate, auCertificate} = this.state;

        if (variant === 'start') {
            this.props.onSubmit({startDate});
            return;
        }

        if (variant === 'update') {
            this.props.onSubmit({expectedEndDate, auCertificate});
            return;
        }

        this.props.onSubmit({
            expectedEndDate,
            auCertificate: auCertificate || 'unchanged',
        });
    };

    private renderDateField(
        id: string,
        label: React.ReactNode,
        value: string,
        onChange: (value: string) => void,
        help?: React.ReactNode,
        error?: string,
    ): React.ReactNode {
        const bounds = this.getDateBounds();
        return (
            <label className='sickleave-modal__field' htmlFor={id}>
                <span className='sickleave-modal__label'>{label}</span>
                <input
                    id={id}
                    className='sickleave-modal__date'
                    type='date'
                    value={value}
                    min={bounds.minDate}
                    max={bounds.maxDate}
                    onChange={(event) => onChange(event.target.value)}
                    disabled={this.props.submitting}
                    required={true}
                />
                {help && <span className='sickleave-modal__help'>{help}</span>}
                {error && <span className='sickleave-modal__error'>{error}</span>}
            </label>
        );
    }

    private renderSelectField(
        id: string,
        label: React.ReactNode,
        value: string,
        options: Array<{value: string; label: React.ReactNode}>,
        onChange: (value: string) => void,
        optional = false,
        error?: string,
    ): React.ReactNode {
        return (
            <label className='sickleave-modal__field' htmlFor={id}>
                <span className='sickleave-modal__label'>{label}</span>
                <select
                    id={id}
                    className='sickleave-modal__select'
                    value={value}
                    onChange={(event) => onChange(event.target.value)}
                    disabled={this.props.submitting}
                    required={!optional}
                >
                    {!optional && (
                        <option value='' disabled={true} hidden={true}>
                            {'—'}
                        </option>
                    )}
                    {optional && (
                        <option value='unchanged'>
                            <FormattedMessage
                                id='dialog.au.unchanged'
                                defaultMessage='Unchanged'
                            />
                        </option>
                    )}
                    {options.map((option) => (
                        <option key={option.value} value={option.value}>
                            {option.label}
                        </option>
                    ))}
                </select>
                {error && <span className='sickleave-modal__error'>{error}</span>}
            </label>
        );
    }

    private renderForm(): React.ReactNode {
        const {variant, fieldErrors, submitting} = this.props;
        const {startDate, expectedEndDate, auCertificate} = this.state;
        const yesNoOptions = [
            {value: 'yes', label: <FormattedMessage id='dialog.au.yes' defaultMessage='Yes'/>},
            {value: 'no', label: <FormattedMessage id='dialog.au.no' defaultMessage='No'/>},
        ];

        return (
            <form className='sickleave-modal__form' onSubmit={this.handleSubmit}>
                {variant === 'start' && this.renderDateField(
                    'start_date',
                    <FormattedMessage id='dialog.a.start_date' defaultMessage='First sick day'/>,
                    startDate,
                    (value) => this.setState({startDate: value}),
                    <FormattedMessage id='dialog.a.start_date_help' defaultMessage='Cannot be in the future.'/>,
                    fieldErrors.start_date,
                )}

                {variant === 'update' && (
                    <>
                        {this.renderDateField(
                            'expected_end_date',
                            <FormattedMessage id='dialog.b.expected_end' defaultMessage='Expected return date'/>,
                            expectedEndDate,
                            (value) => this.setState({expectedEndDate: value}),
                            <FormattedMessage id='dialog.b.expected_end_help' defaultMessage='Must be on or after your first sick day.'/>,
                            fieldErrors.expected_end_date,
                        )}
                        {this.renderSelectField(
                            'au_certificate',
                            <FormattedMessage id='dialog.b.au_certificate' defaultMessage='Medical certificate (AU)'/>,
                            auCertificate,
                            yesNoOptions,
                            (value) => this.setState({auCertificate: value}),
                            false,
                            fieldErrors.au_certificate,
                        )}
                    </>
                )}

                {variant === 'extend' && (
                    <>
                        {this.renderDateField(
                            'expected_end_date',
                            <FormattedMessage id='dialog.c.expected_end' defaultMessage='New expected return date'/>,
                            expectedEndDate,
                            (value) => this.setState({expectedEndDate: value}),
                            <FormattedMessage id='dialog.c.expected_end_help' defaultMessage='Must be after your current expected return date.'/>,
                            fieldErrors.expected_end_date,
                        )}
                        {this.renderSelectField(
                            'au_certificate',
                            <FormattedMessage id='dialog.c.au_certificate' defaultMessage='Medical certificate (AU)'/>,
                            auCertificate,
                            yesNoOptions,
                            (value) => this.setState({auCertificate: value}),
                            true,
                            fieldErrors.au_certificate,
                        )}
                    </>
                )}

                {this.props.generalError && (
                    <div className='sickleave-modal__error sickleave-modal__error--general'>
                        {this.props.generalError}
                    </div>
                )}

                <div className='sickleave-modal__actions'>
                    <button
                        type='button'
                        className='sickleave-modal__button'
                        onClick={this.props.onClose}
                        disabled={submitting}
                    >
                        <FormattedMessage id='dialog.cancel' defaultMessage='Cancel'/>
                    </button>
                    <button
                        type='submit'
                        className='sickleave-modal__button sickleave-modal__button--primary'
                        disabled={submitting}
                    >
                        <FormattedMessage id='dialog.submit' defaultMessage='Submit'/>
                    </button>
                </div>
            </form>
        );
    }

    private titleForVariant(): React.ReactNode {
        switch (this.props.variant) {
        case 'start':
            return <FormattedMessage id='dialog.a.title' defaultMessage='Report sick leave'/>;
        case 'update':
            return <FormattedMessage id='dialog.b.title' defaultMessage='Update sick leave'/>;
        case 'extend':
            return <FormattedMessage id='dialog.c.title' defaultMessage='Extend sick leave'/>;
        default:
            return null;
        }
    }

    private introForVariant(): React.ReactNode {
        switch (this.props.variant) {
        case 'start':
            return <FormattedMessage id='dialog.a.intro' defaultMessage='Enter the first day you were unable to work due to illness.'/>;
        case 'update':
            return <FormattedMessage id='dialog.b.intro' defaultMessage='Provide your expected return date and whether you have a medical certificate (AU).'/>;
        case 'extend':
            return <FormattedMessage id='dialog.c.intro' defaultMessage='Enter a new expected return date if your sick leave continues.'/>;
        default:
            return null;
        }
    }

    render(): React.ReactNode {
        if (!this.props.visible || !this.props.variant) {
            return null;
        }

        const style = getStyle(this.props.theme);

        return (
            <div
                className='sickleave-modal__overlay'
                style={style.overlay}
                onClick={this.props.onClose}
            >
                <div
                    className='sickleave-modal__panel'
                    style={style.panel}
                    onClick={(event) => event.stopPropagation()}
                    role='dialog'
                    aria-modal='true'
                >
                    <h2 className='sickleave-modal__title'>{this.titleForVariant()}</h2>
                    <p className='sickleave-modal__intro'>{this.introForVariant()}</p>
                    {this.renderForm()}
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
