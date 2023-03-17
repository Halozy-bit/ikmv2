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
                <!-- <div class="d-lg-flex flex-column justify-content-center">
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
                </div> -->
                <div class="d-lg-flex justify-content-lg-between flex-wrap justify-content-sm-center">
                <?php
        
                    if (isset($_REQUEST['page'])) {
                        $page = (int)$_REQUEST['page'];
                    } else {
                        $page = 1;
                    }
                    if ($page < 1) $page = 1;
                    $str = file_get_contents('http://192.168.1.67:8082/catalog/'.$page);
                    $json = json_decode($str, true);
                    foreach ($json['catalog'] as $key => $value) {
                        echo "
                        <div class='d-flex flex-column col-lg-3 ps-3'>
                            <div class='pt-3 pe-3'><a href='product.php?id={$value['id']}'><img src='{$value['thumbnail']}' class='img-fluid' alt='Goole Drive Image' style='border-radius: 5%;' height='180'></a></div>
                            <div style='background-color:#FEFBE8;' class='col-7 mt-2 badge rounded-pill text-dark'><a href='' style='text-decoration:none; color:black;'>{$value['kategori']}</a></div>
                            <div class=''><a class='fw-bold fs-4 text-decoration-none text-dark' href='product.php?id={$value['id']}'>{$value['nama']}</a></div>
                            <div><a href='perprofile.php?id={$value['owner']}' style='text-decoration: none; color:black;'>{$value['owner']}</a></span></button></div>
                        </div>
                        ";
                    }
                    $str = file_get_contents('http://192.168.1.67:8082/page/count');
                    $json = json_decode($str, true);
                ?>
                </div>
                <div class="container d-flex justify-content-center pt-3">
                <nav aria-label="Page navigation example">
                    <ul class="pagination">
                        <?php
                            $maxPage = $json['total'];
                            $next = $maxPage - $page;
                            $href = "http://localhost/ikmv2/frontend/view/catalog.php?page=";
                            
                            if ($page == 1) {
                                echo "
                                <li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page) . "'>$page</a></li>";
                                if ($next > 0) {
                                    echo "<li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page+1) . "'>Next</a></li>
                                    <li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page+1) . "'>Next</a></li>";
                                }
                            } elseif($next > 1) {
                                echo "
                                <li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page-1) . "'>Previous</a></li>
                                <li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page-1) . "'>" . ($page-1) . "</a></li>
                                <li class='page-item'><a class='text-dark page-link' href='{$href}$page'>$page</a></li>
                                <li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page+1) . "'>" . ($page+1) . "</a></li>
                                <li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page+1) . "'>Next</a></li>
                            ";
                            } else {
                                echo "
                                <li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page-1) . "'>Previous</a></li>
                                <li class='page-item'><a class='text-dark page-link' href='{$href}" . $page . "'>".($page-1)."</a></li>
                                <li class='page-item'><a class='text-dark page-link' href='{$href}" . ($page+1) . "'>$page</a></li>
                            ";
                            }  
                        ?>
                    </ul>
                </nav>
                </div>
            </div>
        </div>
    </div>