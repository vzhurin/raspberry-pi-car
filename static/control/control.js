$(document).ready(function () {
    speedCoefficient = 1
    timeout = 100

    setInterval(function (){
        bfVal = Number($('#control-back-forward').val());
        lrVal = Number($('#control-left-right').val());

        console.log({bf: bfVal, lr: lrVal})

        lVal = rVal = (-bfVal / 100)

        if (lrVal < 0) {
            lVal = lVal + lrVal

            if (Math.abs(lVal) > 1) {
                if (lVal < 0) {
                    lVal = -1
                }

                if (lVal > 1) {
                    lVal = 1
                }
            }
        }

        if (lrVal > 0) {
            rVal = rVal - lrVal

            if (Math.abs(rVal) > 1) {
                if (rVal < 0) {
                    rVal = -1
                }

                if (rVal > 1) {
                    rVal = 1
                }
            }
        }

        console.log({l: lVal, r: rVal})

        $.ajax("http://192.168.1.35:8070/api/move", {
            data: JSON.stringify({"leftDutyCycle": lVal, "rightDutyCycle": rVal}),
            contentType: 'application/json',
            type: 'POST'
        })

    }, timeout)

    $('#control-stop').on('click', function (){
        $('#control-back-forward').val(0);
        $('#control-left-right').val(0);
    })
})