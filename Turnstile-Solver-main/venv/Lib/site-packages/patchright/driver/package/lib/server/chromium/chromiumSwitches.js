"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.chromiumSwitches = void 0;
/**
 * Copyright 2017 Google Inc. All rights reserved.
 * Modifications copyright (c) Microsoft Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// No dependencies as it is used from the Electron loader.

const chromiumSwitches = exports.chromiumSwitches = ['--disable-field-trial-config',
// https://source.chromium.org/chromium/chromium/src/+/main:testing/variations/README.md
'--disable-background-networking', '--disable-background-timer-throttling', '--disable-backgrounding-occluded-windows',
// Avoids surprises like main request not being intercepted during page.goBack().
'--disable-breakpad',
// Avoids unneeded network activity after startup.
'--no-default-browser-check', '--disable-dev-shm-usage', '--disable-hang-monitor', '--disable-prompt-on-repost', '--disable-renderer-backgrounding', '--force-color-profile=srgb', '--no-first-run', '--password-store=basic', '--use-mock-keychain',
// See https://chromium-review.googlesource.com/c/chromium/src/+/2436773
'--no-service-autorun', '--export-tagged-pdf',
// https://chromium-review.googlesource.com/c/chromium/src/+/4853540
'--disable-search-engine-choice-screen', '--disable-blink-features=AutomationControlled'];