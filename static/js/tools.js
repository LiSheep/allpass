
function gotoURL(url) {
    token = window.localStorage.getItem("token");
    console.log(token)
    $.ajax({
        type: 'get',
        url: url,
        headers:{"Authentication": token},
        success: function(data) {
            $(data).replaceAll("html")
        }
    });
}