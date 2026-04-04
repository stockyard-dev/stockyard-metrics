package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Metrics</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center;flex-wrap:wrap}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.metrics-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:.5rem}
.metric{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem .9rem;transition:border-color .2s}
.metric:hover{border-color:var(--leather)}
.metric-name{font-size:.72rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.3rem}
.metric-value{font-size:1.4rem;font-weight:700}
.metric-unit{font-size:.6rem;color:var(--cm);margin-left:.2rem}
.metric-meta{font-size:.5rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.4rem;flex-wrap:wrap;align-items:center}
.metric-actions{display:flex;gap:.2rem;margin-top:.3rem}
.type-badge{font-size:.45rem;padding:.1rem .3rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.type-badge.counter{border-color:var(--green);color:var(--green)}.type-badge.gauge{border-color:var(--blue);color:var(--blue)}.type-badge.histogram{border-color:var(--gold);color:var(--gold)}.type-badge.timer{border-color:var(--rust);color:var(--rust)}
.tag{font-size:.45rem;padding:.1rem .25rem;background:var(--bg3);color:var(--cm)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.5rem;padding:.15rem .35rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:92vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
@media(max-width:600px){.metrics-grid{grid-template-columns:1fr}.row2{grid-template-columns:1fr}}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> METRICS</h1><button class="btn btn-p" onclick="openForm()">+ Record</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search metrics..." oninput="render()">
<select class="filter-sel" id="type-filter" onchange="render()"><option value="">All Types</option><option value="counter">Counter</option><option value="gauge">Gauge</option><option value="histogram">Histogram</option><option value="timer">Timer</option></select>
</div>
<div class="metrics-grid" id="metrics"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',metrics=[],editId=null;

async function load(){var r=await fetch(A+'/metrics').then(function(r){return r.json()});metrics=r.metrics||[];renderStats();render();}

function renderStats(){
var total=metrics.length;
var sources={};metrics.forEach(function(m){if(m.source)sources[m.source]=true});
var types={};metrics.forEach(function(m){types[m.type]=(types[m.type]||0)+1});
document.getElementById('stats').innerHTML=[
{l:'Metrics',v:total},{l:'Sources',v:Object.keys(sources).length},{l:'Types',v:Object.keys(types).length}
].map(function(x){return '<div class="st"><div class="st-v">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var tf=document.getElementById('type-filter').value;
var f=metrics;
if(tf)f=f.filter(function(m){return m.type===tf});
if(q)f=f.filter(function(m){return(m.name||'').toLowerCase().includes(q)||(m.source||'').toLowerCase().includes(q)||(m.tags||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('metrics').innerHTML='<div class="empty">No metrics recorded. POST to /api/metrics to start.</div>';return;}
var h='';f.forEach(function(m){
h+='<div class="metric">';
h+='<div class="metric-name">'+esc(m.name)+'</div>';
h+='<div><span class="metric-value">'+fmtVal(m.value)+'</span>';
if(m.unit)h+='<span class="metric-unit">'+esc(m.unit)+'</span>';
h+='</div>';
h+='<div class="metric-meta">';
h+='<span class="type-badge '+(m.type||'gauge')+'">'+esc(m.type||'gauge')+'</span>';
if(m.source)h+='<span>'+esc(m.source)+'</span>';
if(m.tags){m.tags.split(',').forEach(function(t){t=t.trim();if(t)h+='<span class="tag">'+esc(t)+'</span>';});}
h+='<span>'+ft(m.created_at)+'</span>';
h+='</div>';
h+='<div class="metric-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+m.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+m.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
});
document.getElementById('metrics').innerHTML=h;
}

function fmtVal(v){if(v>=1000000)return(v/1000000).toFixed(1)+'M';if(v>=1000)return(v/1000).toFixed(1)+'k';return v;}

async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/metrics/'+id,{method:'DELETE'});load();}

function formHTML(metric){
var i=metric||{name:'',type:'gauge',value:0,unit:'',source:'',tags:''};
var isEdit=!!metric;
var h='<h2>'+(isEdit?'EDIT METRIC':'RECORD METRIC')+'</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'" placeholder="e.g. cpu_usage"></div>';
h+='<div class="row2"><div class="fr"><label>Type</label><select id="f-type">';
['counter','gauge','histogram','timer'].forEach(function(t){h+='<option value="'+t+'"'+(i.type===t?' selected':'')+'>'+t.charAt(0).toUpperCase()+t.slice(1)+'</option>';});
h+='</select></div><div class="fr"><label>Value</label><input id="f-value" type="number" value="'+i.value+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Unit</label><input id="f-unit" value="'+esc(i.unit)+'" placeholder="e.g. %, ms, bytes"></div>';
h+='<div class="fr"><label>Source</label><input id="f-source" value="'+esc(i.source)+'" placeholder="e.g. web-01"></div></div>';
h+='<div class="fr"><label>Tags</label><input id="f-tags" value="'+esc(i.tags)+'" placeholder="comma separated"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Record')+'</button></div>';
return h;
}

function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-name').focus();}
function openEdit(id){var m=null;for(var j=0;j<metrics.length;j++){if(metrics[j].id===id){m=metrics[j];break;}}if(!m)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(m);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}

async function submit(){
var name=document.getElementById('f-name').value.trim();
if(!name){alert('Name is required');return;}
var body={name:name,type:document.getElementById('f-type').value,value:parseInt(document.getElementById('f-value').value)||0,unit:document.getElementById('f-unit').value.trim(),source:document.getElementById('f-source').value.trim(),tags:document.getElementById('f-tags').value.trim()};
if(editId){await fetch(A+'/metrics/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/metrics',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}

function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
