    <!-- hero -->
    <div id="hero">
        <div class="container mt-3">
            <p class="fw-bold" style="color: #B54708;">Blog Kami</p>
            <p class="fs-3 fw-bold">Produk IKM UMKM Kecamatan Pujon</p>
            <p class="col-lg-6">katalog dibawah merupakan produk-produk dari IKM dan UMKM sebagai peran aktif pihak terkait untuk mendorong produktivitas serta kemandirian IKM dan UMKM melalui literasi digital</p>
        </div>
    </div>
    <!-- catalog -->
    <div id="catalog">
        <div class="container">
            <div class="container">
                <div class="d-lg-flex flex-column justify-content-center">
                    <form class="d-flex" role="search">
                        <input class="form-control me-2" type="search" width="240" placeholder="Search" aria-label="Search">
                    </form>
                    <div class="container text-center">
                        <button type="button" class="btn btn-warning my-3 fw-bold">Unduh Katalog</button>
                        <p class="fw-bold">Kategori Produk</p>
                        <button type="button" class="btn">Lihat Semua</button>
                        <button type="button" class="btn">Makanan Kering</button>
                        <button type="button" class="btn">Makanan Basah</button>
                        <button type="button" class="btn">Minuman</button>
                        <button type="button" class="btn">Aksesoris</button>
                    </div>
                </div>
                <div class="d-lg-flex justify-content-lg-between flex-wrap justify-content-sm-center">
                <?php
                    $str = file_get_contents('http://192.168.1.67:8082/catalog/1');
                    $json = json_decode($str, true);
                    // var_dump($json);
                    foreach ($json['catalog'] as $key => $value) {
                        echo "
                        <div class='d-flex flex-column col-lg-3 ps-3'>
                            <div class='pt-3 pe-3'><img src='https://lh4.googleusercontent.com/21JKU08_lRAcreQiSuJx1kI1g1I6_Fnsq0X7FrZLxhxplWaK1VYOVqcLVD7yUnVfGuw=w2400' class='img-fluid' style='border-radius: 5%;' height='180'></div>
                            <div style='background-color:#FEFBE8' class='col-7 mt-2 badge rounded-pill text-dark'>{$value['kategori']}</div>
                            <div class='fw-bold fs-4'>{$value['nama']}</div>
                            <div><a href='index.php' style='text-decoration: none; color:black;'>{$value['owner']}</a></span></button></div>
                        </div>
                        ";
                        // print_r($value);
                    }
                ?>
                </div>
                <div class="container d-flex justify-content-center pt-3">
                <nav aria-label="Page navigation example">
                    <ul class="pagination">
                        <li class="page-item"><a class="text-dark page-link" href="#">Previous</a></li>
                        <li class="page-item"><a class="text-dark page-link" href="#">1</a></li>
                        <li class="page-item"><a class="text-dark page-link" href="#">2</a></li>
                        <li class="page-item"><a class="text-dark page-link" href="#">3</a></li>
                        <li class="page-item"><a class="text-dark page-link" href="#">Next</a></li>
                    </ul>
                </nav>
                </div>
            </div>
        </div>
    </div>