export function isNull(obj) {
    return obj === undefined || obj === null || typeof (obj) === 'undefined'
}


export function isBlank(obj) {
    return obj === undefined || obj === null || typeof (obj) === 'undefined' || obj.trim() === '' || obj.trim() === 'null'
}

export function isNullDefault(obj, defaultObj) {
    return isNull(obj) ? '' : defaultObj;
}

export function getLocalStorage(key) {
    let value = localStorage.getItem(key);
    if (isBlank(value)){
        return "";
    }
    return value.trim();
}

export function getUrlParams(name, str) {
    const reg = new RegExp(`(^|&)${ name}=([^&]*)(&|$)`);
    const r = str.substr(1).match(reg);
    if (r != null) return  decodeURIComponent(r[2]); return null;
}

