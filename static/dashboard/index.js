
const color_palette = [
    '#F94144',
    '#F3722C',
    '#F8961E',
    '#F9844A',
    '#F9C74F',
    '#90BE6D',
    '#43AA8B',
    '#4D908E',
    '#577590',
    '#277DA1'
]

window.onload = function() {
    var rm_chart_inner = {
        type: 'pie',
        data: {
            datasets: [{
                data: {},
                backgroundColor: color_palette,
            }]
        }
    };

    var rh_chart_inner = {
        type: 'pie',
        data: {
            datasets: [{
                data: {},
                backgroundColor: color_palette,
            }]
        }
    };

    var rp_chart_inner = {
        type: 'pie',
        data: {
            datasets: [{
                data: {},
                backgroundColor: color_palette,
            }]
        }
    };

    const rm_ctx = document.getElementById('rm_chart').getContext('2d');
    const rh_ctx = document.getElementById('rh_chart').getContext('2d');
    const rp_ctx = document.getElementById('rp_chart').getContext('2d');

    const rm_chart = new Chart(rm_ctx, rm_chart_inner);
    const rh_chart = new Chart(rh_ctx, rh_chart_inner);
    const rp_chart = new Chart(rp_ctx, rp_chart_inner);

    window.trove_update = function() {
        rm_chart_inner.data.labels = window.trove_request_method_keys;
        rm_chart_inner.data.datasets[0].data = window.trove_request_method_values;
        
        rm_chart.update();

        rh_chart_inner.data.labels = window.trove_request_host_keys;
        rh_chart_inner.data.datasets[0].data = window.trove_request_host_values;

        rh_chart.update();

        rp_chart_inner.data.labels = window.trove_request_path_keys;
        rp_chart_inner.data.datasets[0].data = window.trove_request_path_values;

        rp_chart.update();

        console.log("Updated Dashboard");
    }

    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("dashboard/index.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });
}
