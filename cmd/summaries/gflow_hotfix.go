package summaries

import "io"

const (
	hotfixStartSummaryTpl = `Switched to a new branch 'hotfix/{{ .Version }}'

Summary of actions:
- A new branch 'hotfix/{{ .Version }}' was created, based on '{{ if .BaseRef }}{{ .BaseRef }}{{ else }}master{{ end }}'
- metadata.json version was set to {{ .Version }} and automatically commited for you
- You are now on branch 'hotfix/{{ .Version }}'

Follow-up actions:
- Start committing your hot fixes
- When done, run:

	puppet-modulator gflow hotfix finish

	or

	puppet-modulator gflow hotfix finish -p -q -m "MESSAGE"
`
)

type RenderHotfixStartRendererContext struct {
	Version string
	BaseRef string
}

func RenderHotfixStartSummary(w io.Writer, c RenderHotfixStartRendererContext) error {
	return renderTemplateToWriter(w, hotfixStartSummaryTpl, c)
}
