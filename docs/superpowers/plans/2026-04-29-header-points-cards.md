# Header Points Cards Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在顶部栏显示当日累计、当日消耗、当日结余和总计四个积分小卡片。

**Architecture:** 只改前端 `Tasks.vue`。复用已有 `dailyStats` 和 `user` 数据，不新增接口。新增一个 computed 从今日日期匹配 `dailyStats`，template 增加卡片组，style 增加桌面和手机端样式。

**Tech Stack:** Vue 3 Composition API, Chart.js existing stats data, Vite build.

---

## File Structure

- Modify: `web/src/views/Tasks.vue`
  - Template: 顶部 header 中替换单个 `积分:` 文本为积分卡片组
  - Script: 新增 `todayStats` computed
  - Style: 新增 `.points-cards` 和 `.point-card` 样式

---

### Task 1: Add todayStats computed

**Files:**
- Modify: `web/src/views/Tasks.vue:291-295`

- [ ] **Step 1: Add computed after todayText**

在 `todayText` computed 后面添加：

```javascript
const todayStats = computed(() => {
  const today = new Date().toISOString().slice(0, 10)
  const stat = dailyStats.value.find(s => s.date === today)
  return {
    earn: stat?.earn || 0,
    spend: stat?.spend || 0,
    balance: stat?.balance || 0,
    total: user.value?.points || 0
  }
})
```

- [ ] **Step 2: Verify syntax**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds.

---

### Task 2: Replace header points text with cards

**Files:**
- Modify: `web/src/views/Tasks.vue:5-8`

- [ ] **Step 1: Replace header user-info block**

Replace:

```vue
<div class="user-info">
  <span>{{ user?.username }}</span>
  <span class="points">积分: {{ user?.points || 0 }}</span>
  <button class="logout-btn" @click="logout">退出</button>
</div>
```

With:

```vue
<div class="user-info">
  <span class="username">{{ user?.username }}</span>
  <div class="points-cards">
    <span class="point-card earn">当日累计 {{ todayStats.earn }}</span>
    <span class="point-card spend">当日消耗 {{ todayStats.spend }}</span>
    <span class="point-card balance">当日结余 {{ todayStats.balance }}</span>
    <span class="point-card total">总计 {{ todayStats.total }}</span>
  </div>
  <button class="logout-btn" @click="logout">退出</button>
</div>
```

- [ ] **Step 2: Verify page compiles**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds.

---

### Task 3: Add header card styles

**Files:**
- Modify: `web/src/views/Tasks.vue:704-714`

- [ ] **Step 1: Replace `.points` style with card styles**

Replace the existing `.points` block:

```css
.points {
  color: #34c759;
  font-weight: 600;
}
```

With:

```css
.username {
  font-weight: 500;
  color: #1d1d1f;
}

.points-cards {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  align-items: center;
}

.point-card {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}

.point-card.earn {
  background: #e8f8ee;
  color: #34c759;
}

.point-card.spend {
  background: #ffeceb;
  color: #ff3b30;
}

.point-card.balance {
  background: #fff4e5;
  color: #ff9500;
}

.point-card.total {
  background: #e8f2ff;
  color: #007aff;
}
```

- [ ] **Step 2: Add mobile header styles**

Inside existing `@media (max-width: 600px)` near the top, add:

```css
header {
  align-items: flex-start;
  gap: 8px;
}

.user-info {
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 6px;
}

.points-cards {
  justify-content: flex-end;
  gap: 4px;
}

.point-card {
  padding: 3px 6px;
  font-size: 11px;
}
```

- [ ] **Step 3: Build**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds and dist files update.

---

### Task 4: Verify behavior and restart service

**Files:**
- Modify generated: `web/dist/*`

- [ ] **Step 1: Restart user service**

Run:

```bash
systemctl --user restart srdailytask
systemctl --user status srdailytask --no-pager
```

Expected: service active running.

- [ ] **Step 2: Health check**

Run:

```bash
curl -s http://localhost:18888/health
```

Expected:

```json
{"status":"ok"}
```

- [ ] **Step 3: Manual browser verification**

Open `http://localhost:18888` and confirm:
- 顶部显示用户名、四个积分小卡片、退出按钮
- 今日无数据时卡片显示 0
- 打卡后当日累计和总计变化
- 添加支出后当日消耗、当日结余和总计变化
- 手机端卡片自动换行，不遮挡退出按钮

---

### Task 5: Commit and push

**Files:**
- Modify: `web/src/views/Tasks.vue`
- Modify: `web/dist/*`

- [ ] **Step 1: Check status**

Run:

```bash
cd /home/mgter/srDailyTask && git status --short
```

Expected: only `web/src/views/Tasks.vue`, `web/dist/*`, and plan/spec files from this change are present unless pre-existing unrelated changes remain.

- [ ] **Step 2: Commit relevant files**

Run:

```bash
cd /home/mgter/srDailyTask && git add web/src/views/Tasks.vue web/dist docs/superpowers/plans/2026-04-29-header-points-cards.md && git commit -m "$(cat <<'EOF'
顶部栏新增今日积分小卡片

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

Expected: commit succeeds.

- [ ] **Step 3: Push**

Run:

```bash
cd /home/mgter/srDailyTask && git push
```

Expected: push succeeds.

---

## Self-Review

- Spec coverage: covers four metrics, data source reuse, UI small cards, mobile wrap, verification.
- Placeholder scan: no TBD/TODO/placeholders.
- Type consistency: uses existing `dailyStats`, `user`, Vue computed, and existing build flow.
