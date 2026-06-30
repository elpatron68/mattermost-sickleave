// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import manifest from 'manifest';
import type {Store} from 'redux';

import {getCurrentUserLocale} from 'mattermost-redux/selectors/entities/i18n';

import type {GlobalState} from '@mattermost/types/store';

import {confirmEndCase, endSickLeaveCase, openSickLeaveMenu, openSickLeaveModal, parseSickLeaveCommand} from 'actions/sickleave';
import SickLeaveIcon from 'components/icons/sickleave';
import Root from 'components/root';
import {SET_SICK_LEAVE_CONTEXT} from 'action_types';
import {fetchSickLeaveContext} from 'client';
import reducer from 'reducer';
import {sickLeaveCommandTrigger} from 'selectors';
import type {PluginRegistry} from 'types/mattermost-webapp';

import de from '../i18n/de.json';
import en from '../i18n/en.json';

function getTranslationsForLocale(locale: string) {
    switch (locale) {
    case 'de':
        return de;
    default:
        return en;
    }
}

function translate(state: GlobalState, id: string, defaultMessage: string): string {
    const locale = getCurrentUserLocale(state);
    const translations = getTranslationsForLocale(locale) as Record<string, string>;
    return translations[id] || defaultMessage;
}

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState>): Promise<void> {
        registry.registerReducer(reducer);
        registry.registerRootComponent(Root);
        registry.registerTranslations(getTranslationsForLocale);

        fetchSickLeaveContext().
            then((context) => {
                store.dispatch({type: SET_SICK_LEAVE_CONTEXT, context});
            }).
            catch(() => {
                // Keep default slash command trigger until context loads.
            });

        registry.registerChannelHeaderButtonAction(
            <SickLeaveIcon/>,
            () => {
                const state = store.getState();
                const channelId = state.entities.channels.currentChannelId;
                const teamId = state.entities.channels.channels[channelId]?.team_id || '';
                openSickLeaveMenu(channelId, teamId)(store.dispatch, store.getState);
            },
            translate(store.getState(), 'menu.title', 'Sick leave'),
            translate(store.getState(), 'header.button.tooltip', 'Report sick leave'),
        );

        registry.registerSlashCommandWillBePostedHook((message, args) => {
            const trigger = sickLeaveCommandTrigger(store.getState());
            const subcommand = parseSickLeaveCommand(message, trigger);
            if (subcommand === 'start' || subcommand === 'update' || subcommand === 'extend') {
                openSickLeaveModal(subcommand, args.channel_id, args.team_id)(store.dispatch, store.getState);
                return {};
            }
            if (subcommand === 'end') {
                if (confirmEndCase(store.getState)) {
                    endSickLeaveCase(args.channel_id)(store.dispatch, store.getState);
                }
                return {};
            }
            return {message, args};
        });
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void;
        basename?: string;
    }
}

window.registerPlugin(manifest.id, new Plugin());
