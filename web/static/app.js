import { initClipboardActions } from "./app/clipboard.js";
import { initHistoryBack } from "./app/history.js";
import { initHtmxStatus } from "./app/htmx.js";
import { initPreviewLayout } from "./app/layout.js";
import { initSearchPreview } from "./app/search-preview.js";
import { initSettings } from "./app/settings.js";
import { initTheme } from "./app/theme.js";
import { createToast } from "./app/toast.js";

document.documentElement.dataset.serverigh = "ready";

const queuePreviewOffsetUpdate = initPreviewLayout();
const showToast = createToast();

initTheme();
initSettings();
initSearchPreview();
initHtmxStatus(queuePreviewOffsetUpdate);
initClipboardActions(showToast);
initHistoryBack();
