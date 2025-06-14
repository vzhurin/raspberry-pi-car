$(document).ready(function () {
    setInterval(function (){
        lval = $('#control-left').val() / 100;
        rval = $('#control-right').val() / 100;

        if (lval === 0 && rval === 0) {
            return
        }

        console.log({left: lval, right: rval})

        $.ajax("http://192.168.1.35:8070/api/move", {
            data: JSON.stringify({"rightDutyCycle": rval, "leftDutyCycle": lval}),
            contentType: 'application/json',
            type: 'POST'
        })

    }, 500)

    $('#control-stop').on('click', function (){
        $('#control-left').val(0);
        $('#control-right').val(0);
    })
})