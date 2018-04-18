package profiler

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"go.aoe.com/flamingo/core/csrfPreventionFilter"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// profileController shows information about a requested profile
	profileController struct{}
)

const profileTemplate = `<!doctype html>
<html lang="en">
<head>
	<title>Profile</title>
	<link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet">
	<style>
	body {
		background: #475157;
		margin: 0;
		padding: 25px;
		font-family: 'Roboto', sans-serif;
		font-size: 14px;
	}

	.profiler {
		background: #fff;
		padding: 25px;
		max-width: 1000px;
		margin: 0 auto;
	}

	.profiler-summary {
		list-style: none;
		margin: 0;
		padding: 20px;
		background: #F79223;
		color: #fff;
		font-size: 18px;
	}

	.profiler-entries {
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.profiler-entry {
		overflow: hidden;
		padding: 10px 0 5px 10px;
		margin-bottom: 5px;
		background: rgba(50, 50, 50, 0.04);
		border-bottom: 1px solid rgba(50, 50, 50, 0.05);
	}

	.profiler-entry .fnc {
		color: rgab(0, 0, 0, 0.5);
		margin: 0 0 5px;
	}

	.profiler-entry .file {
		overflow: hidden;
		transition: all 0.3s ease-in-out;
		max-height: 0;
	}

	.profiler-subentries {
		margin-top: 10px;
		padding-left: 0px;
	}

	.profiler-entry .duration {
		float: right;
		font-weight: bold;
		font-size: 16px;
	}

	.profiler-entry .duration-relative {
		float: right;
		width: 100px;
		border: 1px solid #ccc;
		margin: 5px 15px;
		clear: both;
	}

	.profiler-entry .duration-relative .offset {
		display: block;
		height: 12px;
		width: 0;
		float: left;
	}

	.profiler-entry .duration-relative .inner {
		background: #F79223;
		display: block;
		height: 12px;
		width: 0;
		float: left;
	}

	.profiler-entry .msg {
		margin: 0;
		font-size: 16px;
		display: inline;
	}

	.profiler-entry .fnc {
		color: rgba(0, 0, 0, 0.5);
	}

	.profiler-entry .fnc.has-file {
		cursor: pointer;
	}

	.profiler-entry .fnc.has-file .icon {
		display: inline-block;
		position: relative;
		top: 0.1em;
		width: 0.75em;
		height: 0.75em;
		border-radius: 2px;
		margin-right: 0.25em;
		border: 1px solid rgba(0, 0, 0, 0.5);
		font-style: normal;
		pointer-events: none;
	}

	.profiler-entry .fnc.has-file .icon:after {
		content: "\002B";
		position: absolute;
		top: -2px;
		left: 2px;
		font-size: 13px;
	}

	.profiler-entry.is-open > .fnc.has-file .icon:after {
		content: "\2013";
		left: 1px;
	}

	.profiler-entry .file-meta,
	.profiler-entry .file-wrap {
		overflow: hidden;
		transition: all 0.3s ease-in-out;
		max-height: 0;
	}

	.profiler-entry.is-open > .file-meta,
	.profiler-entry.is-open > .file-wrap {
		max-height: 500px;
		overflow-y: auto;
	}

	.profiler-entry .file-meta {
		background: #F79223;
		color: white;
		padding: 0 20px 0 10px;
		clear: both;
	}

	.profiler-entry.is-open > .file-meta {
		padding: 5px 20px 5px 10px
	}

	.profiler-entry .file-hint {
		margin: 0 0 10px;
		padding: 10px;
		font-size: 12px;
		line-height: 1.15em;
		background: #fff;
	}

	.profiler-subentries {
		margin-top: 10px;
		padding-left: 0px;
	}
	</style>
</head>
<body>
<div class="profiler">
	<header class="profiler-header">
		<h1>Profile</h1>
		<ul class="profiler-summary">
			<li class="duration-total" data-duration="{{printf "%d" .Duration }}" data-start="{{printf "%d" .Start.UnixNano }}">Time: {{.Duration}}</li>
			<li>Start: {{.Start}}</li>
		</ul>
	</header>
	<div class="profiler-content">
		<h3>Collected Data</h3>
		<ul class="profiler-entries">
		{{ range $data := .Data }}
			<li class="profiler-entry">{{ $data }}</li>
		{{ end }}
		</ul>

		<h3>Profile</h3>
		<ul class="profiler-entries">
		{{ range $entry := .Childs }}
			{{ template "entry" $entry }}
		{{ end }}
		</ul>
	</div>
</div>
<script>
	document.addEventListener('click', function(e) {
		if(e.target.classList.contains('fnc')) {
		e.target.parentNode.classList.toggle('is-open');
		}
	});

	var totalDuration = document.querySelector('.duration-total').dataset.duration;
	var totalStart = document.querySelector('.duration-total').dataset.start;

	Array.from(document.querySelectorAll('.duration-relative')).forEach(addRelativeDuration);
	function addRelativeDuration(element) {
		var duration = element.dataset.duration;
		var relativeDuration = Math.max(1, Math.min(Math.round(100 / totalDuration * duration), 100));
		element.querySelector('.inner').style.width = relativeDuration + '%';

		var start = element.dataset.start;
		var offsetDuration = Math.max(0, Math.min(Math.round(100 / totalDuration * (start - totalStart)), 99));
		element.querySelector('.offset').style.width = offsetDuration + '%';
	}
</script>
</body>
</html>

{{ define "entry" }}
<li class="profiler-entry">
	<span class="duration-relative" data-duration="{{printf "%d" .Duration }}" data-start="{{printf "%d" .Start.UnixNano }}"><i class="offset"></i><i class="inner"></i></span>
	<span class="duration">{{ .Duration }}</span>
	<h3 class="msg">{{ .Msg }}</h3>
	{{- if .Link }}<span class="fnc"><a href="{{ .Link }}">{{ .Link }}</a></span>{{- end}}
	<span class="fnc {{if and .Startpos .Endpos}}has-file{{end}}"><i class="icon"></i>{{ .Fnc }}</span>

	{{if and .Startpos .Endpos}}
		<div class="file-meta">
			<span class="file-path">{{.File}}</span>
			<span class="file-lines">Line {{ .Startpos }} - {{ .Endpos }}</span>
		</div>
		<div class="file-wrap">
			<pre class="file-hint">{{ .Filehint }}</pre>
		</div>
	{{ end }}

	{{if .Childs}}
		<ul class="profiler-subentries">
			{{ range $entry := .Childs }}
			{{ template "entry" $entry }}
			{{ end }}
		</ul>
	{{ end }}
</li>
{{ end }}
`

// Get Response for Debug Info
func (dc *profileController) Get(ctx web.Context) web.Response {
	t, err := template.New("tpl").Parse(profileTemplate)
	if err != nil {
		panic(err)
	}
	var body = new(bytes.Buffer)

	profile, ok := profilestorage.Load(ctx.MustParam1("profile"))
	if !ok {
		return &web.ContentResponse{
			ContentType: "text/html; charset=utf-8",
			BasicResponse: web.BasicResponse{
				Status: http.StatusNotFound,
			},
		}
	}
	t.ExecuteTemplate(body, "tpl", profile)

	return &web.ContentResponse{
		ContentType: "text/html; charset=utf-8",
		BasicResponse: web.BasicResponse{
			Status: http.StatusOK,
		},
		Body: body,
	}
}

// Post saves offline profiling events
func (dc *profileController) Post(ctx web.Context) web.Response {
	dur, _ := strconv.ParseFloat(ctx.MustForm1("duration"), 64)

	profile, ok := profilestorage.Load(ctx.MustParam1("profile"))
	if !ok {
		return &web.ContentResponse{
			ContentType: "text/html; charset=utf-8",
			BasicResponse: web.BasicResponse{
				Status: http.StatusNotFound,
			},
		}
	}

	profile.ProfileOffline(ctx.MustForm1("key"), ctx.MustForm1("message"), time.Duration(dur*1000*1000))

	return &web.JSONResponse{}
}

func (dc *profileController) CheckOption(option router.ControllerOption) bool {
	return option == csrfPreventionFilter.Ignore
}
