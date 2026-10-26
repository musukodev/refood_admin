// Logic Halaman Status & Instruksi Pembayaran ReFood

let currentOrderId = '';
let pollInterval = null;

function formatRupiah(number) {
    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
        maximumFractionDigits: 0
    }).format(number || 0);
}

document.addEventListener('DOMContentLoaded', () => {
    const params = new URLSearchParams(window.location.search);
    currentOrderId = params.get('order_id') || '';

    if (!currentOrderId) {
        document.getElementById('orderIdText').textContent = 'Order ID tidak ditemukan';
        alert('Order ID tidak ditemukan pada URL.');
        return;
    }

    document.getElementById('orderIdText').textContent = currentOrderId;
    checkStatusManual();

    // Polling setiap 4 detik
    pollInterval = setInterval(pollStatus, 4000);
});

async function pollStatus() {
    if (!currentOrderId) return;
    try {
        const res = await fetch(`/api/checkout/status?order_id=${encodeURIComponent(currentOrderId)}`);
        if (!res.ok) return;
        const data = await res.json();
        if (data.order) {
            renderOrder(data.order, data.simulation_enabled);
        }
    } catch (err) {
        console.error('Polling error:', err);
    }
}

async function checkStatusManual() {
    if (!currentOrderId) return;
    try {
        const res = await fetch(`/api/checkout/status?order_id=${encodeURIComponent(currentOrderId)}`);
        const data = await res.json();
        if (!res.ok) {
            alert(data.error || 'Gagal memuat status order.');
            return;
        }
        if (data.order) {
            renderOrder(data.order, data.simulation_enabled);
        }
    } catch (err) {
        alert('Gagal menghubungi server: ' + err.message);
    }
}

function renderOrder(order, simulationEnabled) {
    document.getElementById('orderIdText').textContent = order.id;
    document.getElementById('orderAmountText').textContent = formatRupiah(order.amount);
    document.getElementById('orderFoodName').textContent = order.food_name || '-';
    document.getElementById('orderStore').textContent = order.store || '-';
    document.getElementById('orderQty').textContent = `${order.qty || 1} Porsi`;
    document.getElementById('orderCustomerName').textContent = `${order.customer_name || '-'} (${order.customer_phone || '-'})`;
    document.getElementById('orderMethod').textContent = formatMethodName(order.method);

    const badge = document.getElementById('statusBadge');
    const badgeText = document.getElementById('statusText');
    const badgeDot = document.getElementById('statusDot');
    const successBox = document.getElementById('successBox');
    const pollIndicator = document.getElementById('pollIndicator');
    const contentArea = document.getElementById('paymentContentArea');
    const devTools = document.getElementById('devToolsBox');

    if (simulationEnabled === false) {
        devTools.style.display = 'none';
    }

    if (order.status === 'SUCCESS') {
        badge.className = 'status-badge success';
        badgeText.textContent = 'LUNAS / BERHASIL';
        badgeDot.style.display = 'none';

        successBox.style.display = 'block';
        pollIndicator.style.display = 'none';
        contentArea.style.display = 'none';

        document.getElementById('pickupCode').textContent = `${order.id.slice(-6)}-RF`;

        if (pollInterval) {
            clearInterval(pollInterval);
            pollInterval = null;
        }
        return;
    }

    if (order.status === 'FAILED') {
        badge.className = 'status-badge failed';
        badgeText.textContent = 'PEMBAYARAN GAGAL';
        badgeDot.style.display = 'none';
        pollIndicator.style.display = 'none';
        if (pollInterval) clearInterval(pollInterval);
        return;
    }

    // Status Pending
    badge.className = 'status-badge pending';
    badgeText.textContent = 'MENUNGGU PEMBAYARAN';
    badgeDot.style.display = 'inline-block';
    successBox.style.display = 'none';
    contentArea.style.display = 'block';

    renderPaymentInstructions(order, contentArea);
}

function renderPaymentInstructions(order, container) {
    if (order.qr_string || order.qr_image || order.method === 'qris') {
        const qrSrc = order.qr_image || (order.qr_string ? `https://api.qrserver.com/v1/create-qr-code/?size=260x260&data=${encodeURIComponent(order.qr_string)}` : '');
        container.innerHTML = `
            <div class="qris-display-box">
                <span style="display: inline-block; padding: 4px 12px; background: #FFE4E6; color: #E11D48; border-radius: 9999px; font-size: 11px; font-weight: 800; margin-bottom: 8px;">
                    📱 QRIS STANDAR NASIONAL
                </span>
                <p style="font-size: 13px; color: #4B5563;">Scan QRIS di bawah menggunakan BCA, Mandiri, ShopeePay, DANA, GoPay, atau OVO:</p>

                ${qrSrc ? `
                    <div class="qris-img-wrap">
                        <img src="${qrSrc}" alt="QRIS Barcode">
                    </div>
                ` : `
                    <div style="padding: 30px; color: #6B7280; font-size: 13px;">QRIS sedang dibuat oleh gateway...</div>
                `}

                <p style="font-size: 11px; color: #9CA3AF;">Nominal yang dibayar harus persis: <strong>${formatRupiah(order.amount)}</strong></p>
            </div>
        `;
        return;
    }

    if (order.va_number || order.method.startsWith('va_')) {
        const vaNum = order.va_number || 'Sedang memuat...';
        container.innerHTML = `
            <div class="va-display-box">
                <span style="display: inline-block; padding: 4px 12px; background: rgba(255, 255, 255, 0.15); color: #FDE68A; border-radius: 9999px; font-size: 11px; font-weight: 800;">
                    🏦 VIRTUAL ACCOUNT ${formatMethodName(order.method).toUpperCase()}
                </span>
                <p style="font-size: 12px; color: #94A3B8; margin-top: 8px;">Transfer ke nomor Virtual Account berikut:</p>

                <div class="va-number" id="vaCodeText">${escapeHtml(vaNum)}</div>

                <button type="button" class="btn-copy-va" onclick="copyVA('${escapeHtml(vaNum)}')">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
                    <span id="copyBtnText">Salin Nomor VA</span>
                </button>

                <p style="font-size: 11px; color: #94A3B8; margin-top: 14px;">
                    Dapat dibayar melalui ATM, Mobile Banking, atau Internet Banking sesuai bank Anda.
                </p>
            </div>
        `;
        return;
    }

    // Default POP (Duitku invoice redirect)
    if (order.payment_url) {
        container.innerHTML = `
            <div class="pop-display-box">
                <h4 style="font-size: 15px; font-weight: 800; color: #1E40AF; margin-bottom: 6px;">
                    Checkout Duitku POP Siap
                </h4>
                <p style="font-size: 13px; color: #3B82F6; margin-bottom: 16px;">
                    Klik tombol di bawah untuk membuka halaman pembayaran resmi Duitku:
                </p>
                <a href="${escapeHtml(order.payment_url)}" target="_blank" class="btn-open-pop">
                    <span>Lanjutkan Pembayaran di Duitku</span>
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.4"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
                </a>
            </div>
        `;
        return;
    }

    container.innerHTML = `
        <div style="text-align: center; padding: 24px; color: #6B7280; font-size: 13px;">
            Menunggu konfirmasi transaksi dari gateway Duitku...
        </div>
    `;
}

function formatMethodName(method) {
    if (!method) return 'Duitku';
    const m = method.toLowerCase();
    if (m === 'pop') return 'Duitku POP (All-in-One)';
    if (m === 'qris') return 'QRIS (ShopeePay / DANA)';
    if (m === 'va_bca') return 'BCA Virtual Account';
    if (m === 'va_mandiri') return 'Mandiri Virtual Account';
    if (m === 'va_bni') return 'BNI Virtual Account';
    return method.toUpperCase();
}

function copyVA(text) {
    if (!text || text.includes('...')) return;
    navigator.clipboard.writeText(text).then(() => {
        const btnText = document.getElementById('copyBtnText');
        if (btnText) {
            btnText.textContent = 'Tersalin!';
            setTimeout(() => { btnText.textContent = 'Salin Nomor VA'; }, 2000);
        }
    }).catch(() => {
        alert('Nomor VA: ' + text);
    });
}

async function simulateSuccess() {
    if (!currentOrderId) return;
    const btn = document.getElementById('simulateSuccessBtn');
    btn.disabled = true;
    btn.textContent = 'Memproses Simulasi...';

    try {
        const res = await fetch(`/api/simulate-payment?order_id=${encodeURIComponent(currentOrderId)}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ order_id: currentOrderId })
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Simulasi gagal.');

        // Sukses
        if (data.order) {
            renderOrder(data.order, true);
        } else {
            await checkStatusManual();
        }
    } catch (err) {
        alert('Error simulasi: ' + err.message);
    } finally {
        btn.disabled = false;
        btn.innerHTML = `
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>Simulasikan Pembayaran Berhasil</span>
        `;
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
