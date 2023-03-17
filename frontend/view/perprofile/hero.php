    <?php
    $str = file_get_contents('http://192.168.1.67:8082/umkm/'. $_REQUEST['id']);
    $json = json_decode($str, true);
    $val = $json['umkm']; 
    ?>
    <!-- hero -->
    <div id="hero" class="container">
        <div class="container text-center pt-2">
            <p class="fw-bold" >Detail UMKM</p>
            <h1><?php echo $val['nama'] ?></h1>
            <p><?php echo $val['cerita'][0] ?></p>
        </div>
    </div>
    <div id="content" class="container">
        <div class="container">
            <div class="d-lg-flex flex-row  justify-content-around">
                <div class="d-flex col-lg-6 justify-content-center">
                    <img class="img-fluid" style="height:320px;" src="<?php echo $val['foto_owner']['f_wawancara'] ?>" alt="">
                </div>
                <div class="my-3 col-lg-6 d-flex flex-column justify-content-center">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-tag-fill" viewBox="0 0 16 16">
                        <path d="M2 1a1 1 0 0 0-1 1v4.586a1 1 0 0 0 .293.707l7 7a1 1 0 0 0 1.414 0l4.586-4.586a1 1 0 0 0 0-1.414l-7-7A1 1 0 0 0 6.586 1H2zm4 3.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"/>
                      </svg>
                    <h3>Deskripsi</h3>
                    <p><?php echo $val['deskripsi'] ?></p>
                    <div class="col-lg-4">
                        <button type="button" class="btn btn-warning">Lihat Katalog</button>
                    </div>
                </div>
            </div>
            <div class="d-lg-flex flex-row-reverse justify-content-around">
                <div class="d-flex col-lg-6 justify-content-center">
                    <img class="img-fluid" style="height:320px;" src="<?php echo $val['foto_owner']['f_produksi'] ?>" alt="">
                </div>
                <div class="py-3 d-flex col-lg-6 flex-column justify-content-center">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-tag-fill" viewBox="0 0 16 16">
                        <path d="M2 1a1 1 0 0 0-1 1v4.586a1 1 0 0 0 .293.707l7 7a1 1 0 0 0 1.414 0l4.586-4.586a1 1 0 0 0 0-1.414l-7-7A1 1 0 0 0 6.586 1H2zm4 3.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"/>
                      </svg>
                    <h3>Legalitas</h3>
                    <p><?php 
                    $mp = $val['marketplace'];
                    $sm = $val['social_media'];
                    echo "Nama {$val['nama']} Badan hukum {$val['badan_hukum']} Branding {$val['branding']} Marketplace {$mp['tokopedia']}, {$mp['shopee']}, {$mp['tiktok_shop']}, {$mp['web']}, Sosial Media {$sm['youtube']}, {$sm['tiktok']}, {$sm['instagram']}, {$sm['facebook']}, ";
                    ?></p>
                    <div class="col-lg-4">
                        <button type="button" class="btn btn-warning">Lihat Katalog</button>
                    </div>
                </div>
            </div>
        </div>
    </div>