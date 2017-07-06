
function encode(data, pass, salt, resultcb) {
    Crypto_scrypt(pass, salt, 16, 1, 1, 32, function(result) {
        data = data + salt;
        var key = [];
        for (var i = 0; i < result.length; i ++) {
            key.push(result[i]);
        }
        var dataBytes = aesjs.util.convertStringToBytes(data);
        var aesCtr = new aesjs.ModeOfOperation.ctr(key, new aesjs.Counter(5));
        var encryptedBytes = aesCtr.encrypt(dataBytes);
        resultcb(btoa(encryptedBytes));
        return;
    });
}

function decode(data, pass, salt, resultcb) {
    Crypto_scrypt(pass, salt, 16, 1, 1, 32, function(result) {
        var key = [];
        for (var i = 0; i < result.length; i++) {
            key.push(result[i]);
        }
        var aesCtr = new aesjs.ModeOfOperation.ctr(key, new aesjs.Counter(5));
        data = atob(data);
        data = data.split(',');
        for (var i = 0; i < data.length; i++) {
            data[i] =  parseInt(data[i]);
        }
        var decryptedBytes = aesCtr.decrypt(data);
        var decodeStr = aesjs.util.convertBytesToString(decryptedBytes);
        var salt_len = salt.length;
        decodeStr = decodeStr.slice(0, -salt_len);
        resultcb(decodeStr)
        return;
    });
}

function showSecret(s) {
	var secret = window.localStorage.getItem("secret");
	var salt = window.localStorage.getItem("salt");
	if (secret == "" || secret == undefined || secret == null) {
		$("#setSecret").modal('show');
		return;
	}
	if (s.decoded == undefined || !s.decoded) {
		decode(s.Secret, secret, salt,
		function(result) {
			s.Secret = result;
			s.decoded = true;
		});
	} else {
		encode(s.Secret, secret, salt,
		function(result) {
			s.Secret = result;
			s.decoded = false;
		});
	}
}


function getPasswordSuccess(vm, res) {
	for (var i = 0; i < res.length; i++) {
		res[i].decoded = false;
		res[i].id = i;
	}
	vm.pass = res;
}

function updateAllPassword(oldSecret, oldSalt, newSecret, newSalt, callback) {
    var token = window.localStorage.getItem("token");
    $.ajax({
        type: "get",
        dataType: "json",
        url: "/password/listAPI",
        headers: {
            Authentication: token
        },
        success: (data) => {
            var results = [];
            var funs = [];
            if (data.length > 0) {
                for (var index in data) {
                    d = data[index];
                    funs.push(new Promise( (resolve) => {
                        var id = d.Id;
                        var password = d.Secret;
                        new Promise( (resolve) => {
                            decode(password, oldSecret, oldSalt, resolve);
                        }).
                        then( (ps) => {
                            encode(ps, newSecret, newSalt, (d) => {
                                results.push({
                                    Id: id,
                                    Secret: d
                                });
                                resolve();
                            });
                        });
                    }));
                }// for
                Promise.all(funs).then(function() {
                    $.ajax({
                        type: 'post',
                        dataType: 'json',
                        url: "/password/updateAll",
                        data: {
                            data: JSON.stringify(results)
                        },
                        success: callback
                    });
                });

            } else {
                callback({success: true});
            }
        }
    });
}

function flush_password(vm,cb) {
    var token = window.localStorage.getItem("token");
	$.ajax({
		type: "get",
		dataType: "json",
		url: "/password/listAPI",
		headers: {
			Authentication: token
		},
		success: function(data) {
			getPasswordSuccess(vm, data);
			for (var i = 0; i < data.length; i++) {
				dataCache[i] = {
					id: data[i].id,
					Site: data[i].Site,
					Username: data[i].Username,
					me: data[i]
				};
			}
			$(".ui.search").search({
				source: dataCache,
				searchFields: ['Site'],
				fields: {
					title: 'Site'
				}
			});
		}
	});
	$(".ui.search").on('change',
        function() {
            var results = $('.ui.search').search('get results');
            var data = [];
            for (var i = 0; i < results.length; i++) {
                data[i] = dataCache[results[i].id].me;
            }
    });
}
