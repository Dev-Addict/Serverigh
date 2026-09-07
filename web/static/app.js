(() => {
	document.documentElement.dataset.serverigh = 'ready';

	const updatePreviewOffset = () => {
		const topBar = document.querySelector('.top-bar');
		const breadcrumbs = document.querySelector('.breadcrumbs');
		const footer = document.querySelector('.status-row');
		const previewChrome = document.querySelector('.preview-chrome');
		const topBarHeight = topBar?.offsetHeight || 0;
		const breadcrumbsHeight = breadcrumbs?.offsetHeight || 0;
		const footerHeight = footer?.offsetHeight || 0;
		const previewChromeHeight = previewChrome?.offsetHeight || 0;
		const offset = topBarHeight + breadcrumbsHeight;

		document.documentElement.style.setProperty(
			'--serverigh-top-bar-height',
			`${topBarHeight}px`,
		);

		document.documentElement.style.setProperty(
			'--serverigh-footer-height',
			`${footerHeight}px`,
		);

		document.documentElement.style.setProperty(
			'--serverigh-preview-chrome-height',
			`${previewChromeHeight}px`,
		);

		document.documentElement.style.setProperty(
			'--serverigh-preview-offset',
			`${offset}px`,
		);
	};

	const queuePreviewOffsetUpdate = () => {
		window.requestAnimationFrame(updatePreviewOffset);
	};

	queuePreviewOffsetUpdate();
	window.addEventListener('resize', queuePreviewOffsetUpdate);
	window.addEventListener('load', queuePreviewOffsetUpdate);

	if ('ResizeObserver' in window) {
		const observer = new ResizeObserver(queuePreviewOffsetUpdate);
		const observeChrome = () => {
			document
				.querySelectorAll('.top-bar, .breadcrumbs, .status-row')
				.forEach((node) => {
					observer.observe(node);
				});
			document.querySelectorAll('.preview-chrome').forEach((node) => {
				observer.observe(node);
			});
		};

		observeChrome();
		document.addEventListener('htmx:afterSettle', observeChrome);
	}

	document.addEventListener('htmx:beforeRequest', () => {
		document.documentElement.dataset.loading = 'true';
	});

	document.addEventListener('htmx:afterRequest', () => {
		document.documentElement.dataset.loading = 'false';
		queuePreviewOffsetUpdate();
	});

	document.addEventListener('htmx:responseError', () => {
		document.documentElement.dataset.loading = 'false';
		queuePreviewOffsetUpdate();
	});

	document.addEventListener('htmx:beforeSwap', (event) => {
		const status = event.detail.xhr.status;
		if (status >= 400 && event.detail.xhr.responseText) {
			event.detail.shouldSwap = true;
			event.detail.isError = false;
		}
	});

	document.addEventListener('htmx:afterSettle', queuePreviewOffsetUpdate);
})();
