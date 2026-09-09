package view

var exactFileIcons = mergeIconMaps(
	configFileIcons,
	toolFileIcons,
	documentFileIcons,
)

var compoundExtensionIcons = mergeIconMaps(
	archiveCompoundIcons,
	javascriptCompoundIcons,
	typescriptCompoundIcons,
	configCompoundIcons,
	documentCompoundIcons,
)

var extensionIcons = mergeIconMaps(
	archiveExtensionIcons,
	audioExtensionIcons,
	codeExtensionIcons,
	configExtensionIcons,
	dataExtensionIcons,
	documentExtensionIcons,
	imageExtensionIcons,
	videoExtensionIcons,
	webExtensionIcons,
)

func mergeIconMaps(maps ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, entries := range maps {
		for extension, icon := range entries {
			merged[extension] = icon
		}
	}

	return merged
}
