package summaries

import "io"

const (
	releaseStartSummaryTpl = `Switched to a new branch 'release/{{ .Version }}'

Summary of actions:
- A new branch 'release/{{ .Version }}' was created, based on '{{ if .BaseRef }}{{ .BaseRef }}{{ else }}develop{{ end }}'
- metadata.json version was set to {{ .Version }} and automatically commited for you
- You are now on branch 'release/{{ .Version }}'

Follow-up actions:
- Start committing last-minute fixes in preparing your release
- When done, run:

	puppet-modulator gflow release finish

	or

	puppet-modulator gflow release finish -p -q -m "MESSAGE"
`
)

type RenderReleaseStartRendererContext struct {
	Version string
	BaseRef string
}

func RenderReleaseStartSummary(w io.Writer, c RenderReleaseStartRendererContext) error {
	return renderTemplateToWriter(w, releaseStartSummaryTpl, c)
}
