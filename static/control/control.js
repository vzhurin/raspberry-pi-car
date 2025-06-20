$(document).ready(function () {
    let carHost = window.location.hostname
    let speedCoefficient = 1
    let timeout = 100

    let lEl = $("#control-left");
    let rEl = $("#control-right");

    setInterval(function (){
        let lVal = (Number(lEl.val()) / -100) * speedCoefficient;
        let rVal = (Number(rEl.val()) / -100) * speedCoefficient;

        if (rVal === 0 && lVal === 0) {
            return
        }

        console.log({l: lVal, r: rVal})

        $.ajax("http://"+carHost+":8070/api/move", {
            data: JSON.stringify({"leftDutyCycle": lVal, "rightDutyCycle": rVal}),
            contentType: "application/json",
            type: "POST"
        })

    }, timeout)

    lEl.on("change", function () {
        lEl.val(0)
    });

    rEl.on("change", function () {
        rEl.val(0)
    });
})