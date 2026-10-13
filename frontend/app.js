// Web Component untuk Sidebar ReFood
class RefoodSidebar extends HTMLElement {
    connectedCallback() {
        const currentPath = window.location.pathname;

        // Tentukan menu mana yang aktif
        const isBeranda = currentPath === '/dashboard' || currentPath.endsWith('dashboard.html');
        const isVerifikasi = currentPath === '/verifikasi' || currentPath === '/tinjau' || currentPath.endsWith('verifikasi.html') || currentPath.endsWith('tinjau.html');
        const isPengguna = currentPath === '/pengguna' || currentPath.endsWith('pengguna.html');
        const isKategori = currentPath === '/kategori' || currentPath.endsWith('kategori.html') || currentPath === '/tambah-kategori' || currentPath.endsWith('tambah-kategori.html');
        const isPelaporan = currentPath === '/pelaporan' || currentPath.endsWith('pelaporan.html');
        const isPengaturan = currentPath === '/pengaturan' || currentPath.endsWith('pengaturan.html');

        // Deteksi apakah sedang dibuka lewat server Golang (http://localhost:8080) atau double-click file lokal (file://)
        const isFileProtocol = window.location.protocol === 'file:';

        const linkDashboard = isFileProtocol ? './dashboard.html' : '/dashboard';
        const linkVerifikasi = isFileProtocol ? './verifikasi.html' : '/verifikasi';
        const linkPengguna = isFileProtocol ? './pengguna.html' : '/pengguna';
        const linkKategori = isFileProtocol ? './kategori.html' : '/kategori';
        const linkPelaporan = isFileProtocol ? './coming-soon.html' : '/pelaporan';
        const linkPengaturan = isFileProtocol ? './coming-soon.html' : '/pengaturan';
        const linkLogout = isFileProtocol ? './index.html' : '/';

        this.innerHTML = `
            <aside class="sidebar">
                <div class="sidebar-header">
                    <div class="logo-icon bg-white">
                        <img src="./assets/logo.jpg" alt="ReFood" style="width: 100%; height: 100%; object-fit: contain; border-radius: 4px;">
                    </div>
                    <div class="logo-text-group">
                        <h2>ReFood</h2>
                        <p>Save Food, Share Good</p>
                    </div>
                </div>

                <nav class="sidebar-nav">
                    <a href="${linkDashboard}" class="nav-item ${isBeranda ? 'active' : ''}">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                        Beranda
                    </a>
                    <a href="${linkVerifikasi}" class="nav-item ${isVerifikasi ? 'active' : ''}">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
                        Verifikasi Toko
                    </a>
                    <a href="${linkPengguna}" class="nav-item ${isPengguna ? 'active' : ''}">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
                        Pengguna
                    </a>
                    <a href="${linkKategori}" class="nav-item ${isKategori ? 'active' : ''}">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="12 2 2 7 12 12 22 7 12 2"></polygon><polyline points="2 17 12 22 22 17"></polyline><polyline points="2 12 12 17 22 12"></polyline></svg>
                        Kelola Kategori Produk
                    </a>
                    <a href="${linkPelaporan}" class="nav-item ${isPelaporan ? 'active' : ''}">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="20" x2="18" y2="10"></line><line x1="12" y1="20" x2="12" y2="4"></line><line x1="6" y1="20" x2="6" y2="14"></line></svg>
                        Pelaporan
                    </a>
                    <a href="${linkPengaturan}" class="nav-item ${isPengaturan ? 'active' : ''}">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
                        Akun & Pengaturan
                    </a>
                </nav>

                <div class="user-profile">
                    <div class="profile-info">
                        <div class="avatar-small">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                        </div>
                        <span>Budi Santoso</span>
                    </div>
                    <a href="${linkLogout}" class="logout-btn" title="Keluar">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
                    </a>
                </div>
            </aside>
        `;
    }
}
customElements.define('refood-sidebar', RefoodSidebar);

// Global Function untuk Tanggal Topbar
function initTopbarDate() {
    const currentDateEl = document.getElementById('currentDate');
    if (currentDateEl) {
        const dateOptions = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
        currentDateEl.textContent = new Date().toLocaleDateString('id-ID', dateOptions);
    }
}

// Inisialisasi Berdasarkan Halaman
document.addEventListener('DOMContentLoaded', () => {
    initTopbarDate();

    // 1. Logika Halaman Login
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        const togglePassword = document.getElementById('togglePassword');
        const passwordInput = document.getElementById('password');
        const errorMessage = document.getElementById('errorMessage');
        const successMessage = document.getElementById('successMessage');
        const submitBtn = document.getElementById('submitBtn');

        if (togglePassword && passwordInput) {
            togglePassword.addEventListener('click', () => {
                const type = passwordInput.getAttribute('type') === 'password' ? 'text' : 'password';
                passwordInput.setAttribute('type', type);
                togglePassword.innerHTML = type === 'text' 
                    ? '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>'
                    : '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>';
            });
        }

        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            errorMessage.style.display = 'none';
            successMessage.style.display = 'none';

            const username = document.getElementById('username').value;
            const password = passwordInput.value;

            if (!username || !password) {
                errorMessage.textContent = 'Harap isi Username/Email dan Kata Sandi.';
                errorMessage.style.display = 'block';
                return;
            }

            submitBtn.disabled = true;
            submitBtn.querySelector('span').textContent = 'Memproses...';

            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });
                const data = await response.json();

                if (response.ok) {
                    successMessage.textContent = 'Login berhasil! Mengalihkan...';
                    successMessage.style.display = 'block';
                    if (data.token) localStorage.setItem('refood_token', data.token);
                    setTimeout(() => {
                        window.location.href = window.location.protocol === 'file:' ? './dashboard.html' : '/dashboard';
                    }, 1200);
                } else {
                    errorMessage.textContent = data.error || 'Terjadi kesalahan saat login.';
                    errorMessage.style.display = 'block';
                }
            } catch (err) {
                errorMessage.textContent = 'Gagal terhubung ke server backend Golang.';
                errorMessage.style.display = 'block';
            } finally {
                submitBtn.disabled = false;
                submitBtn.querySelector('span').textContent = 'Masuk Sekarang';
            }
        });
    }

    // 2. Logika Dashboard
    if (document.getElementById('transactionChart')) {
        initDashboard();
    }

    // 3. Logika Verifikasi Toko
    if (document.getElementById('tabMenunggu')) {
        initVerifikasi();
    }

    // 4. Logika Manajemen Pengguna
    if (document.getElementById('usersList')) {
        initPengguna();
    }

    // 5. Logika Kelola Kategori Produk
    if (document.getElementById('kategoriList')) {
        initKategori();
    }

    // 6. Logika Tambah Kategori Baru
    if (document.getElementById('addSymbolsGrid')) {
        initTambahKategori();
    }
});

// Implementasi Dashboard
async function initDashboard() {
    const iconMap = {
        'croissant': '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"/></svg>',
        'coffee': '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8h1a4 4 0 0 1 0 8h-1"></path><path d="M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z"></path><line x1="6" y1="1" x2="6" y2="4"></line><line x1="10" y1="1" x2="10" y2="4"></line><line x1="14" y1="1" x2="14" y2="4"></line></svg>',
        'utensils': '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"></path><path d="M7 2v20"></path><path d="M21 15V2v0a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3zm-3 7v-7"></path></svg>'
    };

    let dashboardData = {
        stats: {
            activeUsers: 1428,
            savedPortions: 215,
            newPartners: 8,
            incidents: 2
        },
        transactions: {
            completed: 142,
            surplusValue: "Rp 3.850.000",
            chartData: [15, 18, 22, 20, 25, 30, 38, 42, 36, 48, 65, 88, 125, 95, 60, 45, 38, 30, 42, 68, 85, 98, 105, 110],
            chartLabels: ["06.00", "", "", "", "", "", "09.00", "", "", "", "", "", "12.00", "", "", "", "", "", "15.00", "", "", "", "", "17.30 (Now)"]
        },
        pendingPartners: [
            { name: "Rotikita Bakery Mer...", type: "Surplus roti sore • Jakp", icon: "croissant" },
            { name: "Kopi Pagi Senopati", type: "Pastry & cold brew • Jal...", icon: "coffee" },
            { name: "Dapoer Nusantara M...", type: "Makanan siap saji • Jaks...", icon: "utensils" }
        ]
    };

    if (window.location.protocol !== 'file:') {
        try {
            const response = await fetch('/api/dashboard');
            if (response.ok) {
                dashboardData = await response.json();
            }
        } catch (err) {
            console.log("Menggunakan data fallback:", err);
        }
    }

    const elUsers = document.getElementById('valActiveUsers');
    const elPortions = document.getElementById('valSavedPortions');
    const elPartners = document.getElementById('valNewPartners');
    const elIncidents = document.getElementById('valIncidents');
    const elCompleted = document.getElementById('valCompletedTx');
    const elSurplus = document.getElementById('valSurplusValue');

    if (elUsers) elUsers.textContent = dashboardData.stats.activeUsers.toLocaleString('id-ID');
    if (elPortions) elPortions.textContent = dashboardData.stats.savedPortions;
    if (elPartners) elPartners.textContent = dashboardData.stats.newPartners;
    if (elIncidents) elIncidents.textContent = dashboardData.stats.incidents;
    if (elCompleted) elCompleted.textContent = dashboardData.transactions.completed;
    if (elSurplus) elSurplus.textContent = dashboardData.transactions.surplusValue;

    const partnersList = document.getElementById('partnersList');
    if (partnersList) {
        partnersList.innerHTML = '';
        dashboardData.pendingPartners.forEach(partner => {
            const isFile = window.location.protocol === 'file:';
            const linkTinjau = isFile ? './tinjau.html' : '/tinjau';
            const html = `
                <div class="partner-item">
                    <div class="partner-info">
                        <div class="partner-icon">${iconMap[partner.icon] || iconMap['coffee']}</div>
                        <div class="partner-text">
                            <h4>${partner.name}</h4>
                            <p>${partner.type}</p>
                        </div>
                    </div>
                    <button class="btn-tinjau" onclick="window.location.href='${linkTinjau}'">Tinjau</button>
                </div>
            `;
            partnersList.insertAdjacentHTML('beforeend', html);
        });
    }

    renderTransactionChart(dashboardData.transactions);
}

function renderTransactionChart(transData) {
    const canvas = document.getElementById('transactionChart');
    if (!canvas || typeof Chart === 'undefined') return;

    const ctx = canvas.getContext('2d');
    const gradient = ctx.createLinearGradient(0, 0, 0, 260);
    gradient.addColorStop(0, 'rgba(56, 110, 85, 0.40)');   
    gradient.addColorStop(1, 'rgba(56, 110, 85, 0.0)');

    if (window.refoodChartInstance) {
        window.refoodChartInstance.destroy();
    }

    window.refoodChartInstance = new Chart(ctx, {
        type: 'line',
        data: {
            labels: transData.chartLabels,
            datasets: [{
                label: 'Volume Transaksi',
                data: transData.chartData,
                borderColor: '#386E55',
                backgroundColor: gradient,
                borderWidth: 2.5,
                fill: true,
                tension: 0.42,
                pointRadius: 0,
                pointHoverRadius: 6,
                pointHoverBackgroundColor: '#386E55',
                pointHoverBorderColor: '#FFFFFF',
                pointHoverBorderWidth: 2
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
                intersect: false,
                mode: 'index',
            },
            plugins: { 
                legend: { display: false },
                tooltip: {
                    backgroundColor: '#1A1A1A',
                    titleFont: { size: 11, family: 'Plus Jakarta Sans' },
                    bodyFont: { size: 12, family: 'Plus Jakarta Sans', weight: 'bold' },
                    padding: 10,
                    cornerRadius: 8,
                    displayColors: false,
                    callbacks: {
                        label: function(context) {
                            return `${context.parsed.y} Transaksi`;
                        }
                    }
                }
            },
            scales: {
                x: { 
                    grid: { display: false, drawBorder: false }, 
                    ticks: { 
                        font: { size: 11, family: 'Plus Jakarta Sans', weight: '600' }, 
                        color: '#757575',
                        autoSkip: false
                    } 
                },
                y: { 
                    display: false, 
                    min: 0,
                    suggestedMax: 140
                }
            }
        }
    });
}

// Implementasi Verifikasi Toko
async function initVerifikasi() {
    const tabMenunggu = document.getElementById('tabMenunggu');
    const tabTerverifikasi = document.getElementById('tabTerverifikasi');
    const storesList = document.getElementById('storesList');
    let currentTab = 'menunggu';
    let data = { menunggu: [], terverifikasi: [] };

    try {
        const response = await fetch('/api/verifikasi');
        data = await response.json();
        document.getElementById('countMenunggu').textContent = data.menunggu.length;
        document.getElementById('countTerverifikasi').textContent = data.terverifikasi.length;
        renderList();
    } catch (err) {
        // Fallback offline
        data = {
            menunggu: [
                { id: 1, name: "Artisan Bakery", owner: "Sri Wahyuni", location: "Batam Center", time: "Diajukan hari ini, 09:15 WIB", status: "Menunggu", image: "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop" },
                { id: 2, name: "Artisan Bakery", owner: "Sri Wahyuni", location: "Batam Center", time: "Diajukan hari ini, 09:15 WIB", status: "Menunggu", image: "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop" },
                { id: 3, name: "Artisan Bakery", owner: "Sri Wahyuni", location: "Batam Center", time: "Diajukan hari ini, 09:15 WIB", status: "Menunggu", image: "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop" }
            ],
            terverifikasi: [
                { id: 101, name: "Artisan Bakery", owner: "Sri Wahyuni", location: "Batam Center", time: "Disetujui hari ini, 09:15 WIB", status: "Terverifikasi", image: "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop" },
                { id: 102, name: "Artisan Bakery", owner: "Sri Wahyuni", location: "Batam Center", time: "Disetujui hari ini, 09:15 WIB", status: "Terverifikasi", image: "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop" }
            ]
        };
        document.getElementById('countMenunggu').textContent = data.menunggu.length;
        document.getElementById('countTerverifikasi').textContent = data.terverifikasi.length;
        renderList();
    }

    tabMenunggu.addEventListener('click', () => {
        currentTab = 'menunggu';
        tabMenunggu.classList.add('active');
        tabTerverifikasi.classList.remove('active');
        document.getElementById('countMenunggu').className = 'tab-badge green';
        document.getElementById('countTerverifikasi').className = 'tab-badge gray';
        renderList();
    });

    tabTerverifikasi.addEventListener('click', () => {
        currentTab = 'terverifikasi';
        tabTerverifikasi.classList.add('active');
        tabMenunggu.classList.remove('active');
        document.getElementById('countMenunggu').className = 'tab-badge gray';
        document.getElementById('countTerverifikasi').className = 'tab-badge green';
        renderList();
    });

    function renderList() {
        storesList.innerHTML = '';
        const currentData = data[currentTab] || [];
        if (currentData.length === 0) {
            storesList.innerHTML = `<div style="text-align:center; padding:40px; color:#757575;">Tidak ada data toko.</div>`;
            return;
        }

        const isFile = window.location.protocol === 'file:';
        const linkTinjau = isFile ? './tinjau.html' : '/tinjau';

        currentData.forEach(store => {
            const isMenunggu = currentTab === 'menunggu';
            const badgeClass = isMenunggu ? 'menunggu' : 'terverifikasi';
            const btnClass = isMenunggu ? 'green' : 'orange';
            const btnText = isMenunggu ? 'Lihat Kelengkapan Data & Dokumen' : 'Lihat Profil & Dokumen';

            const html = `
                <div class="store-card">
                    <div class="store-info">
                        <img src="${store.image}" alt="Store Image" class="store-image">
                        <div class="store-details">
                            <h3>
                                ${store.name}
                                <span class="status-badge ${badgeClass}">
                                    <span class="status-dot dot"></span> ${store.status}
                                </span>
                            </h3>
                            <div class="store-meta bold">${store.owner} &bull; ${store.location}</div>
                            <div class="store-meta">${store.time}</div>
                        </div>
                    </div>
                    <button class="action-btn ${btnClass}" onclick="window.location.href='${linkTinjau}'">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
                        ${btnText}
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg>
                    </button>
                </div>
            `;
            storesList.insertAdjacentHTML('beforeend', html);
        });
    }
}

// Implementasi Manajemen Pengguna
async function initPengguna() {
    const filterSemua = document.getElementById('filterSemua');
    const filterToko = document.getElementById('filterToko');
    const filterDiblokir = document.getElementById('filterDiblokir');
    const usersList = document.getElementById('usersList');
    const paginationInfo = document.getElementById('paginationInfo');
    let currentFilter = 'semua';
    let allData = { total: 0, toko: 0, diblokir: 0, users: [] };

    try {
        const response = await fetch('/api/pengguna');
        allData = await response.json();
    } catch (err) {
        // Fallback data
        let dummy = [];
        for (let i = 0; i < 5; i++) {
            dummy.push({ id: i + 1, name: "Morning Bakery", subtitle: "Budi Santoso (Pemilik) • morningbakery@gmail.com", type: "Usaha", location: "Batam Center, Kepri", status: "Aktif", avatar: "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop", isStore: true });
        }
        for (let i = 5; i < 10; i++) {
            dummy.push({ id: i + 1, name: "Budi Santoso", subtitle: "budisantoso@gmail.com", type: "Individu", location: "Batam Center, Kepri", status: "Aktif", avatar: "BS", isStore: false });
        }
        for (let i = 10; i < 15; i++) {
            dummy.push({ id: i + 1, name: "Budi Santoso", subtitle: "budisantoso@gmail.com", type: "Individu", location: "Batam Center, Kepri", status: "Diblokir", avatar: "BS", isStore: false });
        }
        allData = { total: 1428, toko: 174, diblokir: 5, users: dummy };
    }

    if (filterSemua) filterSemua.querySelector('.chip-count').textContent = allData.total.toLocaleString('id-ID');
    if (filterToko) filterToko.querySelector('.chip-count').textContent = allData.toko.toLocaleString('id-ID');
    if (filterDiblokir) filterDiblokir.querySelector('.chip-count').textContent = allData.diblokir.toLocaleString('id-ID');
    renderList();

    if (filterSemua) filterSemua.addEventListener('click', () => setFilter('semua', filterSemua));
    if (filterToko) filterToko.addEventListener('click', () => setFilter('toko', filterToko));
    if (filterDiblokir) filterDiblokir.addEventListener('click', () => setFilter('diblokir', filterDiblokir));

    function setFilter(filterType, element) {
        currentFilter = filterType;
        [filterSemua, filterToko, filterDiblokir].forEach(btn => btn && btn.classList.remove('active'));
        element.classList.add('active');
        renderList();
    }

    function renderList() {
        if (!usersList) return;
        usersList.innerHTML = '';
        let filteredUsers = allData.users;
        if (currentFilter === 'toko') {
            filteredUsers = allData.users.filter(u => u.isStore);
            if (paginationInfo) paginationInfo.innerHTML = `Menampilkan <b>${filteredUsers.length}</b> dari <b>${allData.toko}</b> Toko Terverifikasi`;
        } else if (currentFilter === 'diblokir') {
            filteredUsers = allData.users.filter(u => u.status === 'Diblokir');
            if (paginationInfo) paginationInfo.innerHTML = `Menampilkan <b>${filteredUsers.length}</b> dari <b>${allData.diblokir}</b> Diblokir`;
        } else {
            if (paginationInfo) paginationInfo.innerHTML = `Menampilkan <b>${filteredUsers.length}</b> dari <b>${allData.total}</b> Semua Pengguna`;
        }

        if (filteredUsers.length === 0) {
            usersList.innerHTML = `<div style="text-align:center; padding:40px; color:#757575;">Tidak ada data pengguna.</div>`;
            return;
        }

        filteredUsers.forEach((user, index) => {
            const avatarHtml = user.isStore ? `<img src="${user.avatar}" alt="${user.name}">` : `${user.avatar}`;
            const iconStore = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>`;
            const iconUser = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>`;
            const iconLocation = `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>`;
            const typeIcon = user.type === 'Usaha' ? iconStore : iconUser;
            const statusClass = user.status.toLowerCase();
            const isSelected = (currentFilter === 'semua' && index === 8) ? 'selected' : '';

            const html = `
                <div class="table-row ${isSelected}">
                    <div class="cell-pengguna">
                        <div class="avatar-wrapper">${avatarHtml}</div>
                        <div class="user-text">
                            <h4>${user.name}</h4>
                            <p>${user.subtitle}</p>
                        </div>
                    </div>
                    <div class="cell-tipe">${typeIcon} ${user.type}</div>
                    <div class="cell-lokasi">${iconLocation} ${user.location}</div>
                    <div class="cell-status">
                        <span class="status-pill ${statusClass}">
                            <span class="status-dot dot"></span> ${user.status}
                        </span>
                    </div>
                </div>
            `;
            usersList.insertAdjacentHTML('beforeend', html);
        });
    }
}

// Logika Tinjau Dokumen (Modal & Lightbox)
function openApproveModal() {
    const approveModal = document.getElementById('approveModal');
    if (approveModal) approveModal.classList.add('active');
}

function openRejectModal() {
    alert('Fitur Tolak: akan menampilkan alasan penolakan.');
}

function closeModals() {
    const approveModal = document.getElementById('approveModal');
    if (approveModal) approveModal.classList.remove('active');
}

function confirmAction(action) {
    if (action === 'approve') {
        alert('Toko berhasil disetujui! Mengalihkan ke halaman verifikasi...');
        window.location.href = window.location.protocol === 'file:' ? './verifikasi.html' : '/verifikasi';
    }
}

function openLightbox(imageUrl) {
    const lightboxImage = document.getElementById('lightboxImage');
    const lightboxOverlay = document.getElementById('lightboxOverlay');
    if (lightboxImage && lightboxOverlay) {
        lightboxImage.src = imageUrl;
        lightboxOverlay.classList.add('active');
    }
}

function closeLightbox() {
    const lightboxOverlay = document.getElementById('lightboxOverlay');
    if (lightboxOverlay) lightboxOverlay.classList.remove('active');
}

// ==========================================================
// LOGIKA HALAMAN KELOLA & TAMBAH KATEGORI PRODUK
// ==========================================================
let kategoriData = [
    { id: 1, name: "Makanan Berat", desc: "Nasi, lauk pauk, makanan utama siap santap higienis dari resto dan katering.", status: true, icon: "🍱" },
    { id: 2, name: "Makanan Ringan", desc: "Kue basah, gorengan, snack kemasan, dan camilan layak konsumsi.", status: true, icon: "🥐" },
    { id: 3, name: "Minuman", desc: "Kopi, teh, jus segar, dan minuman tersegel belum kedaluwarsa.", status: true, icon: "☕" },
    { id: 4, name: "Bahan Makanan", desc: "Sayuran segar, buah utuh, bahan mentah layak olah dari pasar & panen.", status: true, icon: "🥔" },
    { id: 5, name: "Non Konsumsi", desc: "Makanan tidak layak konsumsi manusia untuk kompos & pakan maggot.", status: true, icon: "🥫" }
];

let selectedKategoriId = null;

function initKategori() {
    const listEl = document.getElementById('kategoriList');
    if (!listEl) return;

    renderKategoriList();

    const pickItems = document.querySelectorAll('#editIconPicker .icon-pick-item');
    pickItems.forEach(item => {
        item.addEventListener('click', () => {
            pickItems.forEach(i => {
                i.classList.remove('active');
                const check = i.querySelector('.check-mark');
                if (check) check.remove();
            });
            item.classList.add('active');
            item.insertAdjacentHTML('beforeend', '<span class="check-mark">✓</span>');
        });
    });

    const editDesc = document.getElementById('editKatDesc');
    const editCount = document.getElementById('editCharCount');
    if (editDesc && editCount) {
        editDesc.addEventListener('input', () => {
            editCount.textContent = `${editDesc.value.length} / 300`;
        });
    }
}

function renderKategoriList() {
    const listEl = document.getElementById('kategoriList');
    if (!listEl) return;
    listEl.innerHTML = '';

    kategoriData.forEach(kat => {
        const row = document.createElement('div');
        row.className = 'kategori-row';
        row.innerHTML = `
            <div class="kat-name-cell">
                <div class="kat-icon-box">${kat.icon}</div>
                <span>${kat.name}</span>
            </div>
            <div class="kat-desc-cell">${kat.desc}</div>
            <div class="kat-status-cell">
                <label class="switch-container">
                    <input type="checkbox" ${kat.status ? 'checked' : ''} onchange="toggleKategoriStatus(${kat.id}, this.checked)">
                    <span class="slider"></span>
                    <span class="switch-label">${kat.status ? 'Aktif' : 'Nonaktif'}</span>
                </label>
            </div>
            <div class="kat-action-cell">
                <button class="btn-icon-action" title="Edit" onclick="openEditModal(${kat.id})">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                </button>
                <button class="btn-icon-action delete" title="Hapus" onclick="openDeleteModal(${kat.id})">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                </button>
            </div>
        `;
        listEl.appendChild(row);
    });
}

function toggleKategoriStatus(id, checked) {
    const item = kategoriData.find(k => k.id === id);
    if (item) item.status = checked;
    renderKategoriList();
}

function openEditModal(id) {
    selectedKategoriId = id;
    const item = kategoriData.find(k => k.id === id);
    if (!item) return;

    document.getElementById('editKatName').value = item.name;
    document.getElementById('editKatDesc').value = item.desc;
    document.getElementById('editKatStatus').checked = item.status;
    document.getElementById('editCharCount').textContent = `${item.desc.length} / 300`;

    const modal = document.getElementById('editKategoriModal');
    if (modal) modal.classList.add('active');
}

function closeEditModal() {
    const modal = document.getElementById('editKategoriModal');
    if (modal) modal.classList.remove('active');
}

function saveEditKategori() {
    const item = kategoriData.find(k => k.id === selectedKategoriId);
    if (item) {
        item.name = document.getElementById('editKatName').value;
        item.desc = document.getElementById('editKatDesc').value;
        item.status = document.getElementById('editKatStatus').checked;
        renderKategoriList();
    }
    closeEditModal();
}

function openDeleteModal(id) {
    selectedKategoriId = id;
    const modal = document.getElementById('deleteKategoriModal');
    if (modal) modal.classList.add('active');
}

function closeDeleteModal() {
    const modal = document.getElementById('deleteKategoriModal');
    if (modal) modal.classList.remove('active');
}

function confirmDeleteKategori() {
    kategoriData = kategoriData.filter(k => k.id !== selectedKategoriId);
    renderKategoriList();
    closeDeleteModal();
}

function initTambahKategori() {
    const symbols = document.querySelectorAll('#addSymbolsGrid .symbol-card');
    symbols.forEach(sym => {
        sym.addEventListener('click', () => {
            symbols.forEach(s => s.classList.remove('active'));
            sym.classList.add('active');
        });
    });

    const addDesc = document.getElementById('addKatDesc');
    const addCount = document.getElementById('addCharCount');
    if (addDesc && addCount) {
        addDesc.addEventListener('input', () => {
            addCount.textContent = `${addDesc.value.length} / 300 karakter`;
        });
    }
}

function submitTambahKategori() {
    const name = document.getElementById('addKatName').value;
    if (!name) {
        alert('Harap isi nama kategori produk!');
        return;
    }
    alert(`Kategori "${name}" berhasil disimpan dan diterbitkan!`);
    window.location.href = window.location.protocol === 'file:' ? './kategori.html' : '/kategori';
}
