package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/System9-Software/sysix/internal/collector"
)

func Start(port int) error {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/snapshot", handleSnapshot)
	http.HandleFunc("/api/ports", handlePorts)
	http.HandleFunc("/api/network", handleNetwork)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("sysix web UI running at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashboard))
}

func handleSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot, err := collector.GetSnapshot()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

func handlePorts(w http.ResponseWriter, r *http.Request) {
	ports, err := collector.GetPorts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ports)
}

func handleNetwork(w http.ResponseWriter, r *http.Request) {
	network, err := collector.GetNetwork()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(network)
}

const dashboard = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>sysix — System9</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { background: #090D15; color: #E8F0FF; font-family: monospace; padding: 24px; }
  h1 { color: #4DA8FF; font-size: 1.4rem; margin-bottom: 4px; }
  .subtitle { color: #6B7C99; font-size: 0.85rem; margin-bottom: 24px; }
  .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 16px; }
  .card { background: #1E2D45; border-radius: 8px; padding: 16px; }
  .card h2 { color: #4DA8FF; font-size: 0.9rem; margin-bottom: 12px; text-transform: uppercase; letter-spacing: 0.05em; }
  .row { display: flex; justify-content: space-between; margin-bottom: 6px; }
  .label { color: #6B7C99; font-size: 0.85rem; }
  .value { color: #E8F0FF; font-size: 0.85rem; }
  .good { color: #39D98A; }
  .warn { color: #F5A623; }
  .bad { color: #FF5B5B; }
  table { width: 100%; border-collapse: collapse; font-size: 0.8rem; }
  th { color: #6B7C99; text-align: left; padding: 4px 8px; border-bottom: 1px solid #1E2D45; }
  td { color: #E8F0FF; padding: 4px 8px; }
  tr:hover { background: #090D15; }
</style>
</head>
<body>
<h1>sysix</h1>
<p class="subtitle">System9 Observer — live system monitor</p>
<div class="grid">
  <div class="card">
    <h2>System</h2>
    <div class="row"><span class="label">Host</span><span class="value" id="host">—</span></div>
    <div class="row"><span class="label">OS</span><span class="value" id="os">—</span></div>
    <div class="row"><span class="label">Uptime</span><span class="value" id="uptime">—</span></div>
    <div class="row"><span class="label">CPU</span><span class="value" id="cpu">—</span></div>
    <div class="row"><span class="label">Memory</span><span class="value" id="mem">—</span></div>
    <div class="row"><span class="label">Disk</span><span class="value" id="disk">—</span></div>
  </div>
  <div class="card">
    <h2>Network</h2>
    <div class="row"><span class="label">Sent</span><span class="value" id="sent">—</span></div>
    <div class="row"><span class="label">Received</span><span class="value" id="recv">—</span></div>
    <div class="row"><span class="label">Packets Out</span><span class="value" id="pkts-out">—</span></div>
    <div class="row"><span class="label">Packets In</span><span class="value" id="pkts-in">—</span></div>
  </div>
  <div class="card" style="grid-column: span 2;">
    <h2>Processes</h2>
    <table>
      <thead><tr><th>PID</th><th>Name</th><th>CPU%</th><th>Mem</th></tr></thead>
      <tbody id="procs"></tbody>
    </table>
  </div>
  <div class="card" style="grid-column: span 2;">
    <h2>Ports</h2>
    <table>
      <thead><tr><th>Port</th><th>Type</th><th>Status</th><th>PID</th></tr></thead>
      <tbody id="ports"></tbody>
    </table>
  </div>
</div>
<script>
function statusClass(pct) {
  if (pct >= 90) return 'bad';
  if (pct >= 70) return 'warn';
  return 'good';
}
function formatBytes(b) {
  if (b >= 1e9) return (b/1e9).toFixed(1) + ' GB';
  if (b >= 1e6) return (b/1e6).toFixed(1) + ' MB';
  if (b >= 1e3) return (b/1e3).toFixed(1) + ' KB';
  return b + ' B';
}
async function refresh() {
  try {
    const [snap, net, ports] = await Promise.all([
      fetch('/api/snapshot').then(r => r.json()),
      fetch('/api/network').then(r => r.json()),
      fetch('/api/ports').then(r => r.json()),
    ]);

    document.getElementById('host').textContent = snap.Hostname;
    document.getElementById('os').textContent = snap.OS;
    document.getElementById('uptime').textContent = Math.floor(snap.Uptime / 3600) + ' hours';

    const cpu = document.getElementById('cpu');
    cpu.textContent = snap.CPUPercent.toFixed(1) + '%';
    cpu.className = 'value ' + statusClass(snap.CPUPercent);

    const mem = document.getElementById('mem');
    mem.textContent = snap.MemPercent.toFixed(1) + '% (' + Math.floor(snap.MemUsed/1024/1024) + ' MB / ' + Math.floor(snap.MemTotal/1024/1024) + ' MB)';
    mem.className = 'value ' + statusClass(snap.MemPercent);

    const disk = document.getElementById('disk');
    disk.textContent = snap.DiskPercent.toFixed(1) + '% (' + Math.floor(snap.DiskUsed/1024/1024/1024) + ' GB / ' + Math.floor(snap.DiskTotal/1024/1024/1024) + ' GB)';
    disk.className = 'value ' + statusClass(snap.DiskPercent);

    document.getElementById('sent').textContent = formatBytes(net.BytesSent);
    document.getElementById('recv').textContent = formatBytes(net.BytesRecv);
    document.getElementById('pkts-out').textContent = net.PacketsSent;
    document.getElementById('pkts-in').textContent = net.PacketsRecv;

    const procsBody = document.getElementById('procs');
    procsBody.innerHTML = snap.Processes
      .filter(p => p.MemMB > 1)
      .slice(0, 15)
      .map(p => '<tr><td>'+p.PID+'</td><td>'+p.Name+'</td><td>'+p.CPUPercent.toFixed(1)+'%</td><td>'+p.MemMB.toFixed(0)+'MB</td></tr>')
      .join('');

    const portsBody = document.getElementById('ports');
    portsBody.innerHTML = ports
      .slice(0, 15)
      .map(p => '<tr><td>'+p.Port+'</td><td>'+p.Type+'</td><td>'+p.Status+'</td><td>'+p.PID+'</td></tr>')
      .join('');

  } catch(e) { console.error(e); }
}
refresh();
setInterval(refresh, 2000);
</script>
</body>
</html>`
