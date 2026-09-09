import { initClipboardActions } from "./app/clipboard.js";
import { initHistoryBack } from "./app/history.js";
import { initHtmxStatus } from "./app/htmx.js";
import { initHelpModal } from "./app/help.js";
import { initKeyboardNavigation } from "./app/keyboard.js";
import { initPreviewLayout } from "./app/layout.js";
import { initSearchPreview } from "./app/search-preview.js";
import { initSettings } from "./app/settings.js";
import { initTheme } from "./app/theme.js";
import { createToast } from "./app/toast.js";
import { initVimStatus } from "./app/vim-status.js";
import { initWriteActions } from "./app/write-actions.js";

document.documentElement.dataset.serverigh = "ready";

const queuePreviewOffsetUpdate = initPreviewLayout();
const showToast = createToast();
const setKeyStatus = initVimStatus();

initTheme();
initSettings();
initHelpModal();
initSearchPreview();
initHtmxStatus(queuePreviewOffsetUpdate);
initClipboardActions(showToast);
initHistoryBack();
initKeyboardNavigation(showToast, setKeyStatus);
initWriteActions(showToast);
