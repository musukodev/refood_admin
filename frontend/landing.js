// ReFood Landing Page Logic

let allFoods = [];
let selectedFood = null;

function formatRupiah(number) {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        maximumFractionDigits: 0
    }).format(number);
}

document.addEventListener('DOMContentLoaded', () => {
    loadFoods();
    setupFilters();
    setupModalEvents();
});

async function loadFoods() {
    const container = document.getElementById('foodsContainer');
    try {
        const res = await fetch('/api/foods');
        if (!res.ok) throw new Error('Gagal memuat katalog');
        const data = await res.json();
        allFoods = data.foods || [];
        renderFoods(allFoods);
    } catch (err) {
        console.error('Error fetching foods:', err);
        container.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: #DC2626;">
                <p style="font-weight: 700; margin-bottom: 8px;">Gagal memuat daftar makanan.</p>
                <p style="font-size: 13px; color: #6B7280;">Pastikan server backend Golang sedang berjalan.</p>
                <button onclick="loadFoods()" style="margin-top: 14px; padding: 8px 18px; border-radius: 9999px; background: #386E55; color: #fff; border: none; cursor: pointer; font-weight: 600;">Coba Lagi</button>
            </div>
        `;
    }
}

function renderFoods(foods) {
    const container = document.getElementById('foodsContainer');
    if (!foods || foods.length === 0) {
        container.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: #6B7280;">
                Tidak ada makanan di kategori ini saat ini.
            </div>
        `;
        return;
    }

    container.innerHTML = foods.map(food => {
        const discountPct = Math.round(((food.original_price - food.price) / food.original_price) * 100);
        return `
            <article class="food-card" data-category="${escapeHtml(food.category)}">
                <div class="food-card-thumb">
                    <img src="${escapeHtml(food.image)}" alt="${escapeHtml(food.name)}" onerror="this.src='./assets/food-box.jpg'">
                    <span class="food-discount-tag">HEMAT ${discountPct}%</span>
                    <span class="food-category-tag">${escapeHtml(food.category)}</span>
                </div>
                <div class="food-card-body">
                    <div class="food-store-row">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                            <circle cx="12" cy="10" r="3"></circle>
                        </svg>
                        <span>${escapeHtml(food.store)}</span>
                    </div>
                    <h3 class="food-title">${escapeHtml(food.name)}</h3>
                    <p class="food-desc">${escapeHtml(food.description)}</p>

                    <div class="food-pricing-row">
                        <span class="food-price-current">${formatRupiah(food.price)}</span>
                        <span class="food-price-original">${formatRupiah(food.original_price)}</span>
                    </div>

                    <div class="food-stock-info">
                        <span class="stock-pill">
                            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                            Sisa ${food.stock} porsi
                        </span>
                        <span>Ambil hari ini</span>
                    </div>

                    <button type="button" class="btn-buy-food" onclick="openCheckoutModal('${escapeHtml(food.id)}')">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                            <circle cx="9" cy="21" r="1"></circle>
                            <circle cx="20" cy="21" r="1"></circle>
                            <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path>
                        </svg>
                        <span>Selamatkan Sekarang</span>
                    </button>
                </div>
            </article>
        `;
    }).join('');
}

function setupFilters() {
    const buttons = document.querySelectorAll('#categoryFilter .filter-btn');
    buttons.forEach(btn => {
        btn.addEventListener('click', () => {
            buttons.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            const cat = btn.getAttribute('data-category');
            if (cat === 'all') {
                renderFoods(allFoods);
            } else {
                const filtered = allFoods.filter(f => f.category === cat);
                renderFoods(filtered);
            }
        });
    });
}

function setupModalEvents() {
    const modal = document.getElementById('checkoutModal');
    const closeBtn = document.getElementById('modalCloseBtn');
    const cancelBtn = document.getElementById('modalCancelBtn');
    const qtySelect = document.getElementById('checkoutQty');
    const form = document.getElementById('checkoutForm');

    const closeModal = () => {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    };

    closeBtn.addEventListener('click', closeModal);
    cancelBtn.addEventListener('click', closeModal);

    modal.addEventListener('click', (e) => {
        if (e.target === modal) closeModal();
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && modal.classList.contains('active')) {
            closeModal();
        }
    });

    qtySelect.addEventListener('change', updateModalTotal);

    // Visual radio selection
    const methodOptions = document.querySelectorAll('.method-option');
    methodOptions.forEach(opt => {
        opt.addEventListener('click', () => {
            methodOptions.forEach(o => o.classList.remove('selected'));
            opt.classList.add('selected');
            const radio = opt.querySelector('input[type="radio"]');
            if (radio) radio.checked = true;
        });
    });

    // Handle Form Submit
    form.addEventListener('submit', handleCheckoutSubmit);
}

function openCheckoutModal(foodId) {
    selectedFood = allFoods.find(f => f.id === foodId);
    if (!selectedFood) return;

    document.getElementById('checkoutFoodId').value = selectedFood.id;
    document.getElementById('summaryTitle').textContent = selectedFood.name;
    document.getElementById('summaryStore').textContent = selectedFood.store;
    document.getElementById('summaryPrice').textContent = `${formatRupiah(selectedFood.price)} / porsi`;
    document.getElementById('summaryThumb').src = selectedFood.image;
    document.getElementById('summaryThumb').onerror = () => {
        document.getElementById('summaryThumb').src = './assets/food-box.jpg';
    };

    // Reset qty
    document.getElementById('checkoutQty').value = "1";
    updateModalTotal();

    // Reset error
    const errBox = document.getElementById('checkoutError');
    errBox.style.display = 'none';
    errBox.textContent = '';

    const modal = document.getElementById('checkoutModal');
    modal.classList.add('active');
    document.body.style.overflow = 'hidden';
}

function updateModalTotal() {
    if (!selectedFood) return;
    const qty = parseInt(document.getElementById('checkoutQty').value, 10) || 1;
    const total = selectedFood.price * qty;
    document.getElementById('modalTotalValue').textContent = formatRupiah(total);
}

async function handleCheckoutSubmit(e) {
    e.preventDefault();
    const submitBtn = document.getElementById('modalSubmitBtn');
    const errBox = document.getElementById('checkoutError');
    errBox.style.display = 'none';

    const foodId = document.getElementById('checkoutFoodId').value;
    const qty = parseInt(document.getElementById('checkoutQty').value, 10) || 1;
    const customerName = document.getElementById('customerName').value.trim();
    const customerEmail = document.getElementById('customerEmail').value.trim();
    const customerPhone = document.getElementById('customerPhone').value.trim();
    const paymentMethod = document.querySelector('input[name="payment_method"]:checked')?.value || 'pop';

    submitBtn.disabled = true;
    const originalText = submitBtn.innerHTML;
    submitBtn.innerHTML = `
        <svg class="spin-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="12" y1="2" x2="12" y2="6"></line>
            <line x1="12" y1="18" x2="12" y2="22"></line>
            <line x1="4.93" y1="4.93" x2="7.76" y2="7.76"></line>
            <line x1="16.24" y1="16.24" x2="19.07" y2="19.07"></line>
            <line x1="2" y1="12" x2="6" y2="12"></line>
            <line x1="18" y1="12" x2="22" y2="12"></line>
            <line x1="4.93" y1="19.07" x2="7.76" y2="16.24"></line>
            <line x1="16.24" y1="7.76" x2="19.07" y2="4.93"></line>
        </svg>
        <span>Menghubungkan ke Duitku...</span>
    `;

    try {
        const payload = {
            food_id: foodId,
            qty,
            customer_name: customerName,
            customer_email: customerEmail,
            customer_phone: customerPhone,
            payment_method: paymentMethod
        };

        const response = await fetch('/api/checkout', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Terjadi kesalahan saat memproses checkout.');
        }

        // Sukses membuat transaksi!
        const orderId = data.order_id || data.order?.id;

        // Redirect ke halaman pembayaran terpadu
        window.location.href = `/pembayaran?order_id=${encodeURIComponent(orderId)}`;
    } catch (err) {
        console.error('Checkout error:', err);
        errBox.textContent = err.message || 'Gagal memproses pembayaran.';
        errBox.style.display = 'block';
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;
    }
}

function escapeHtml(text) {
    if (!text) return '';
    return String(text)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}
