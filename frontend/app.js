const API_URL = '/api/items';

let items = [];
let currentFilter = '';

const suggestions = [
    '早安问候 💕',
    '晚安亲吻 😘',
    '一起吃饭 🍽️',
    '牵手散步 👫',
    '拥抱取暖 🤗',
    '看日出 🌅',
    '看日落 🌄',
    '一起看电影 🎬',
    '逛街购物 🛍️',
    '烹饪美食 🍳',
    '互送礼物 🎁',
    '甜蜜合照 📸',
    '说“我爱你” 💌',
    '写情书 ✉️',
    '制造惊喜 🎉',
    '按摩放松 💆',
    '一起运动 🏃',
    '听同一首歌 🎵',
    '互相按摩 💆‍♀️',
    '分享零食 🍪'
];

function formatDate(date) {
    const d = new Date(date);
    const month = d.getMonth() + 1;
    const day = d.getDate();
    return `${month}-${day}`;
}

function formatFullDate(date) {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = (d.getMonth() + 1).toString().padStart(2, '0');
    const day = d.getDate().toString().padStart(2, '0');
    return `${year}-${month}-${day}`;
}

async function loadItems() {
    try {
        let url = API_URL;
        if (currentFilter === 'today') {
            url = API_URL + '?date=' + formatFullDate(new Date());
        } else if (currentFilter) {
            url = API_URL + '?date=' + currentFilter;
        }
        
        const res = await fetch(url);
        items = await res.json();
        renderItems();
    } catch (err) {
        console.error('加载失败:', err);
    }
}

function renderItems() {
    const list = document.getElementById('itemList');
    const emptyState = document.getElementById('emptyState');
    const progressBar = document.getElementById('progressBar');
    const completedCount = document.getElementById('completedCount');
    const totalCount = document.getElementById('totalCount');

    list.innerHTML = '';

    const completed = items.filter(item => item.completed).length;
    const total = items.length;
    const percentage = total > 0 ? (completed / total) * 100 : 0;

    progressBar.style.width = `${percentage}%`;
    completedCount.textContent = completed;
    totalCount.textContent = total;

    if (total === 0) {
        emptyState.classList.remove('hidden');
    } else {
        emptyState.classList.add('hidden');
    }

    items.forEach((item, index) => {
        const li = document.createElement('li');
        li.className = item.completed ? 'completed' : '';
        li.style.animationDelay = `${index * 0.05}s`;
        li.dataset.id = item.id;

        const hasItemDate = !!item.item_date;
        const dateStr = hasItemDate ? formatDate(item.item_date) : formatDate(item.created_at);
        const dateLabel = hasItemDate ? '计划日期' : '创建时间';

        li.innerHTML = `
            <div class="checkbox" onclick="toggleItem(${item.id})"></div>
            <div class="item-content">
                <span class="content">${escapeHtml(item.content)}</span>
                <div class="time-row">
                    <span class="time-label">${dateLabel}</span>
                    <span class="time">${dateStr}</span>
                </div>
            </div>
            <button class="delete-btn" onclick="deleteItem(${item.id})" title="删除">×</button>
        `;
        list.appendChild(li);
    });
}

async function addItem() {
    const input = document.getElementById('newItem');
    const dateInput = document.getElementById('itemDate');
    const content = input.value.trim();

    if (!content) return;

    const body = { content };
    if (dateInput.value) {
        body.item_date = dateInput.value;
    }

    try {
        const res = await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        
        if (!res.ok) {
            const err = await res.json();
            alert(err.error || '添加失败');
            return;
        }
        
        const newItem = await res.json();
        items.unshift(newItem);
        input.value = '';
        dateInput.value = '';
        renderItems();
        
        input.focus();
    } catch (err) {
        console.error('添加失败:', err);
        alert('网络错误，请重试');
    }
}

async function toggleItem(id) {
    const item = items.find(i => i.id === id);
    if (!item) return;

    try {
        const res = await fetch(`${API_URL}/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ completed: !item.completed })
        });
        
        if (!res.ok) {
            const err = await res.json();
            alert(err.error || '更新失败');
            return;
        }
        
        const updated = await res.json();
        item.completed = updated.completed;
        renderItems();
    } catch (err) {
        console.error('更新失败:', err);
        alert('网络错误，请重试');
    }
}

async function deleteItem(id) {
    const li = document.querySelector(`li[data-id="${id}"]`);
    if (li) {
        li.classList.add('removing');
        await new Promise(r => setTimeout(r, 300));
    }

    try {
        const res = await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
        
        if (!res.ok) {
            const err = await res.json();
            alert(err.error || '删除失败');
            if (li) li.classList.remove('removing');
            return;
        }
        
        items = items.filter(i => i.id !== id);
        renderItems();
    } catch (err) {
        console.error('删除失败:', err);
        alert('网络错误，请重试');
        if (li) li.classList.remove('removing');
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

document.getElementById('addBtn').addEventListener('click', addItem);
document.getElementById('newItem').addEventListener('keypress', e => {
    if (e.key === 'Enter') addItem();
});

document.getElementById('suggestBtn').addEventListener('click', function(e) {
    e.stopPropagation();
    const popup = document.getElementById('suggestPopup');
    const list = document.getElementById('suggestList');
    
    if (popup.classList.contains('hidden')) {
        list.innerHTML = suggestions.map(s => 
            `<div class="suggest-item" onclick="selectSuggestion('${s.replace(/'/g, "\\'")}')">${s}</div>`
        ).join('');
        popup.classList.remove('hidden');
    } else {
        popup.classList.add('hidden');
    }
});

document.addEventListener('click', function(e) {
    const popup = document.getElementById('suggestPopup');
    const btn = document.getElementById('suggestBtn');
    if (!popup.contains(e.target) && !btn.contains(e.target)) {
        popup.classList.add('hidden');
    }
});

function selectSuggestion(text) {
    document.getElementById('newItem').value = text;
    document.getElementById('suggestPopup').classList.add('hidden');
    document.getElementById('newItem').focus();
}

document.querySelectorAll('.filter-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        
        const dateValue = btn.dataset.date;
        if (dateValue === 'today') {
            currentFilter = 'today';
        } else {
            currentFilter = '';
        }
        loadItems();
    });
});

const today = new Date();
const year = today.getFullYear();
const month = (today.getMonth() + 1).toString().padStart(2, '0');
const day = today.getDate().toString().padStart(2, '0');
document.getElementById('itemDate').value = `${year}-${month}-${day}`;

loadItems();