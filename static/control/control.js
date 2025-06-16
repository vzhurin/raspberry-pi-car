$(document).ready(function () {
    let speedCoefficient = 0.7
    let timeout = 100

    let lEl = $('#control-left');
    let rEl = $('#control-right');

    setInterval(function (){
        let lVal = (Number(lEl.val()) / -100) * speedCoefficient;
        let rVal = (Number(rEl.val()) / -100) * speedCoefficient;

        console.log({l: lVal, r: rVal})

        $.ajax("http://192.168.1.35:8070/api/move", {
            data: JSON.stringify({"leftDutyCycle": lVal, "rightDutyCycle": rVal}),
            contentType: 'application/json',
            type: 'POST'
        })

    }, timeout)

    lEl.on('change', function () {
        lEl.val(0)
    });

    rEl.on('change', function () {
        rEl.val(0)
    });
})