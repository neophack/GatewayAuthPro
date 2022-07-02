import fetch from 'cross-fetch';
import MD5 from 'crypto-js/md5'

const api = "/";

const post = (body) => ({
    method: 'POST',
    dataType: 'json',
    headers: {
        'Content-Type': 'application/json;charset=UTF-8'
    },
    body: JSON.stringify(body)
});

export function Login(uri, account, password) {
    let url = api + uri;
    if (uri.startsWith('http') || uri.startsWith('//')) {
        url = uri;
    }
    return new Promise((resolve, reject) => {
        var pd = MD5(password).toString().toUpperCase();
        return fetch(url, post({"account": account, "password": pd}))
            .then(res=>res.json())
            .then(res => {
                resolve(res)
            })
            .catch(e => {
                reject({type: 'catch', e: e})
            })
    })
}