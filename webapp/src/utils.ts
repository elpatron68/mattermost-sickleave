// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import manifest from 'manifest';

export function normalizeLocale(locale: string): string {
    if (locale.toLowerCase().startsWith('de')) {
        return 'de';
    }
    return 'en';
}

export function formatDateForLocale(isoDate: string, locale: string): string {
    const parsed = parseISODate(isoDate);
    if (!parsed) {
        return isoDate;
    }

    const intlLocale = normalizeLocale(locale) === 'de' ? 'de-DE' : 'en-US';
    return new Intl.DateTimeFormat(intlLocale, {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
    }).format(parsed);
}

export function getPluginURL(): string {
    if (window.basename) {
        return `${window.basename}/plugins/${manifest.id}`;
    }
    return `/plugins/${manifest.id}`;
}

export function formatISODate(date: Date): string {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
}

export function addDays(date: Date, days: number): Date {
    const copy = new Date(date.getTime());
    copy.setDate(copy.getDate() + days);
    return copy;
}

export function parseISODate(value: string): Date | null {
    const match = (/^(\d{4})-(\d{2})-(\d{2})$/).exec(value);
    if (!match) {
        return null;
    }
    const year = Number(match[1]);
    const month = Number(match[2]);
    const day = Number(match[3]);
    const parsed = new Date(year, month - 1, day);
    if (parsed.getFullYear() !== year || parsed.getMonth() !== month - 1 || parsed.getDate() !== day) {
        return null;
    }
    return parsed;
}
