package gfutils

type HotfixOrReleaseBaseOptions struct {
	Push           bool
	Message        *string
	MessageFile    *string
	NoEditorPrompt bool
}
