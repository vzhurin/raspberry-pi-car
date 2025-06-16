$(document).ready(function () {
    let speedCoefficient = 0.4
    let timeout = 100

    let bfEl = $('#control-back-forward');
    let lrEl = $('#control-left-right');

    setInterval(function (){
        let bfVal = Number(bfEl.val());
        let lrVal = Number(lrEl.val());

        console.log({bf: bfVal, lr: lrVal})

        let lVal= (-bfVal / 100)
        let rVal = (-bfVal / 100)



        lVal = lVal * speedCoefficient
        rVal = rVal * speedCoefficient

        console.log({l: lVal, r: rVal})

        $.ajax("http://192.168.1.35:8070/api/move", {
            data: JSON.stringify({"leftDutyCycle": lVal, "rightDutyCycle": rVal}),
            contentType: 'application/json',
            type: 'POST'
        })

    }, timeout)

    bfEl.on('change', function () {
        bfEl.val(0)
    });

    lrEl.on('change', function () {
        lrEl.val(0)
    });
})