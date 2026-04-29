# Points Heatmap Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在积分统计图下方新增类似 GitHub 的每日结余热力图。

**Architecture:** 只改 `web/src/views/Tasks.vue`，复用现有 `dailyStats` 数据，不新增接口和依赖。前端将每日数据按周分组，渲染为横向周、纵向星期的小方块网格；颜色由 `balance` 决定，正数绿色、负数红色、零灰色。

**Tech Stack:** Vue 3 Composition API, CSS Grid/Flex, existing daily stats API, Vite.

---

## File Structure

- Modify: `web/src/views/Tasks.vue`
  - Template: 在 `chart-wrapper` 下方增加 `heatmap-section`
  - Script: 新增星期/月标签、`heatmapWeeks` computed、颜色 class 和 tooltip 方法
  - Style: 新增热力图布局、月份标签、星期标签、小方块颜色样式

---

### Task 1: Add heatmap data helpers

**Files:**
- Modify: `web/src/views/Tasks.vue` after `todayStats` computed

- [ ] **Step 1: Add constants and helper functions**

Add this after `todayStats` computed:

```javascript
const weekLabels = ['日', '一', '二', '三', '四', '五', '六']

const formatHeatmapTitle = (day) => {
  if (!day) return ''
  return `${day.date}\n累计 +${day.earn}\n消耗 -${day.spend}\n结余 ${day.balance >= 0 ? '+' : ''}${day.balance}`
}

const heatmapClass = (balance) => {
  if (balance >= 30) return 'pos-3'
  if (balance >= 10) return 'pos-2'
  if (balance > 0) return 'pos-1'
  if (balance <= -30) return 'neg-3'
  if (balance <= -10) return 'neg-2'
  if (balance < 0) return 'neg-1'
  return 'zero'
}
```

- [ ] **Step 2: Build check**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds.

---

### Task 2: Add heatmapWeeks computed

**Files:**
- Modify: `web/src/views/Tasks.vue` after helper functions from Task 1

- [ ] **Step 1: Add computed**

Add this after `heatmapClass`:

```javascript
const heatmapWeeks = computed(() => {
  const stats = dailyStats.value.map(s => ({
    date: s.date,
    earn: s.earn,
    spend: s.spend,
    balance: s.balance,
    weekday: new Date(`${s.date}T00:00:00`).getDay(),
    month: Number(s.date.slice(5, 7))
  }))

  if (stats.length === 0) return []

  const weeks = []
  let currentWeek = Array(7).fill(null)

  stats.forEach((day, index) => {
    if (index === 0 && day.weekday > 0) {
      currentWeek = Array(7).fill(null)
    }

    currentWeek[day.weekday] = day

    const isLastDay = index === stats.length - 1
    if (day.weekday === 6 || isLastDay) {
      weeks.push({
        days: currentWeek,
        month: currentWeek.find(Boolean)?.month || day.month
      })
      currentWeek = Array(7).fill(null)
    }
  })

  return weeks
})
```

- [ ] **Step 2: Build check**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds.

---

### Task 3: Add heatmap template below chart

**Files:**
- Modify: `web/src/views/Tasks.vue:207-210`

- [ ] **Step 1: Insert heatmap section after chart-wrapper**

Find:

```vue
<div class="chart-wrapper">
  <div v-if="loadingStats" class="loading">加载中...</div>
  <canvas v-show="!loadingStats" ref="chartCanvas"></canvas>
</div>
```

Replace with:

```vue
<div class="chart-wrapper">
  <div v-if="loadingStats" class="loading">加载中...</div>
  <canvas v-show="!loadingStats" ref="chartCanvas"></canvas>
</div>

<div class="heatmap-section" v-if="!loadingStats && heatmapWeeks.length > 0">
  <div class="heatmap-title">每日结余热力图</div>
  <div class="heatmap-scroll">
    <div class="heatmap-layout">
      <div class="weekday-labels">
        <span v-for="label in weekLabels" :key="label">{{ label }}</span>
      </div>
      <div class="heatmap-weeks">
        <div class="heatmap-week" v-for="(week, weekIndex) in heatmapWeeks" :key="weekIndex">
          <span class="month-label">{{ weekIndex === 0 || week.month !== heatmapWeeks[weekIndex - 1]?.month ? week.month + '月' : '' }}</span>
          <span
            v-for="(day, dayIndex) in week.days"
            :key="dayIndex"
            class="heatmap-cell"
            :class="day ? heatmapClass(day.balance) : 'empty'"
            :title="formatHeatmapTitle(day)"
          ></span>
        </div>
      </div>
    </div>
  </div>
</div>
```

- [ ] **Step 2: Build check**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds.

---

### Task 4: Add heatmap styles

**Files:**
- Modify: `web/src/views/Tasks.vue` near chart styles after `.chart-wrapper canvas`

- [ ] **Step 1: Add CSS**

Add after `.chart-wrapper canvas` block:

```css
.heatmap-section {
  margin-top: 10px;
  border-top: 1px solid #f0f0f0;
  padding-top: 12px;
}

.heatmap-title {
  color: #1d1d1f;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
}

.heatmap-scroll {
  overflow-x: auto;
  padding-bottom: 4px;
}

.heatmap-layout {
  display: flex;
  gap: 6px;
  min-width: max-content;
}

.weekday-labels {
  display: grid;
  grid-template-rows: repeat(7, 10px);
  gap: 3px;
  padding-top: 17px;
  color: #86868b;
  font-size: 9px;
  line-height: 10px;
}

.heatmap-weeks {
  display: flex;
  gap: 3px;
}

.heatmap-week {
  display: grid;
  grid-template-rows: 14px repeat(7, 10px);
  gap: 3px;
}

.month-label {
  color: #86868b;
  font-size: 9px;
  line-height: 12px;
  white-space: nowrap;
}

.heatmap-cell {
  width: 10px;
  height: 10px;
  border-radius: 2px;
  background: #ebedf0;
}

.heatmap-cell.zero,
.heatmap-cell.empty {
  background: #ebedf0;
}

.heatmap-cell.pos-1 { background: #b7ebc6; }
.heatmap-cell.pos-2 { background: #6bd982; }
.heatmap-cell.pos-3 { background: #34c759; }
.heatmap-cell.neg-1 { background: #ffc1bd; }
.heatmap-cell.neg-2 { background: #ff8a83; }
.heatmap-cell.neg-3 { background: #ff3b30; }
```

- [ ] **Step 2: Add mobile refinement**

Inside existing `@media (max-width: 600px)` add:

```css
.heatmap-title {
  font-size: 12px;
}

.heatmap-cell {
  width: 9px;
  height: 9px;
}

.heatmap-week {
  grid-template-rows: 13px repeat(7, 9px);
}

.weekday-labels {
  grid-template-rows: repeat(7, 9px);
}
```

- [ ] **Step 3: Build check**

Run:

```bash
cd /home/mgter/srDailyTask/web && npm run build
```

Expected: build succeeds and dist files update.

---

### Task 5: Verify, restart, commit, push

**Files:**
- Modify: `web/src/views/Tasks.vue`
- Modify: `web/dist/*`
- Create: `docs/superpowers/plans/2026-04-29-points-heatmap.md`
- Already created: `docs/superpowers/specs/2026-04-29-points-heatmap-design.md`

- [ ] **Step 1: Restart service**

Run:

```bash
systemctl --user restart srdailytask
curl -s http://localhost:18888/health
```

Expected:

```json
{"status":"ok"}
```

- [ ] **Step 2: Manual browser check**

Open `http://localhost:18888` and confirm:
- 热力图在积分统计图下方
- 切换天数时热力图同步变化
- 正结余是绿色，负结余是红色，0是灰色
- 手机端可横向滚动，不挤压页面

- [ ] **Step 3: Check status**

Run:

```bash
cd /home/mgter/srDailyTask && git status --short
```

Expected: includes `web/src/views/Tasks.vue`, `web/dist/*`, heatmap spec and plan.

- [ ] **Step 4: Commit and push**

Run:

```bash
cd /home/mgter/srDailyTask && git add web/src/views/Tasks.vue web/dist docs/superpowers/specs/2026-04-29-points-heatmap-design.md docs/superpowers/plans/2026-04-29-points-heatmap.md && git commit -m "$(cat <<'EOF'
积分统计新增每日结余热力图

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)" && git push
```

Expected: commit and push succeed.

---

## Self-Review

- Spec coverage: covers location under chart, GitHub-like arrangement, date range sync, green/red/gray color mapping, tooltip, mobile scroll.
- Placeholder scan: no placeholders.
- Type consistency: uses existing `dailyStats` fields `date`, `earn`, `spend`, `balance`; helper names match template.
