const API_URL = '/api/items';

let items = [];

async function loadItems() {
    try {
        const res = await fetch(API_URL);
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
        li.innerHTML = `
            <div class="checkbox" onclick="toggleItem(${item.id})"></div>
            <span class="content">${escapeHtml(item.content)}</span>
            <button class="delete-btn" onclick="deleteItem(${item.id})" title="删除">×</button>
        `;
        list.appendChild(li);
    });
}

async function addItem() {
    const input = document.getElementById('newItem');
    const content = input.value.trim();

    if (!content) return;

    try {
        const res = await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ content })
        });
        
        if (!res.ok) {
            const err = await res.json();
            alert(err.error || '添加失败');
            return;
        }
        
        const newItem = await res.json();
        items.push(newItem);
        input.value = '';
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

loadItems();
