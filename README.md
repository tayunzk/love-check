# 💕 恋爱打卡清单

记录情侣之间甜蜜时光的 Todo 应用。

## 功能特性

- ✅ 添加打卡项
- ✅ 标记完成/未完成
- ✅ 删除打卡项
- 📊 进度条实时显示完成比例
- 💾 数据持久化存储

## 技术栈

- **前端**：原生 HTML + CSS + JavaScript
- **后端**：Go + Gin
- **部署**：Docker Compose

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/tayun/love-check.git
cd love-check
```

### 2. 启动服务

```bash
docker-compose up -d
```

### 3. 访问应用

浏览器打开：http://localhost:8080

## 项目结构

```
love-check/
├── frontend/          # 前端静态文件
│   ├── index.html
│   ├── app.js
│   └── style.css
├── backend/           # Go 后端服务
│   ├── main.go
│   ├── Dockerfile
│   └── data/          # 数据存储目录
└── docker-compose.yml
```

## API 接口

| 方法   | 路径           | 说明         |
|--------|----------------|--------------|
| GET    | /api/items     | 获取所有项   |
| POST   | /api/items     | 添加新项     |
| PUT    | /api/items/:id | 切换完成状态 |
| DELETE | /api/items/:id | 删除项       |

## License

MIT