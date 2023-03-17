<?php
// $path selalu gunakan slash(/) sebagai awalan
function CallApi($path) {
    $str = file_get_contents('http://192.168.1.67:8082/product'. $path);
    return json_decode($str, true);
}  