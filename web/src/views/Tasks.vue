<template>
  <div class="dashboard-container">
    <div class="today-info">
      <p>今天是 {{ todayText }}</p>
    </div>

    <div class="section profile-section home-profile-section">
      <div class="section-header">
        <h2>个人账户</h2>
      </div>
      <div class="profile-card">
        <div class="profile-top">
          <div class="profile-avatar">
            <img v-if="user?.avatar_url" :src="user.avatar_url" alt="头像" />
            <span v-else>{{ userInitial }}</span>
          </div>
          <div class="profile-details">
            <span class="profile-title">账户信息</span>
            <h3>{{ user?.username }}</h3>
            <p>{{ user?.email || '未设置邮箱' }}</p>
          </div>
          <div class="profile-actions">
            <button class="edit-profile-btn" @click="openProfileModal">编辑资料</button>
            <button class="logout-btn" @click="logout">退出</button>
          </div>
        </div>
        <div class="points-cards profile-points-cards">
          <span class="point-card earn"><span class="point-label">当日累计</span><span class="point-value">{{ todayStats.earn }}</span></span>
          <span class="point-card spend"><span class="point-label">当日消耗</span><span class="point-value">{{ todayStats.spend }}</span></span>
          <span class="point-card balance"><span class="point-label">当日结余</span><span class="point-value">{{ todayStats.balance }}</span></span>
          <span class="point-card total"><span class="point-label">总计</span><span class="point-value">{{ todayStats.total }}</span></span>
          <span class="point-card daily-cost"><span class="point-label">日均消费</span><span class="point-value">{{ formatMoney(longTermSummary.active_daily_cost) }}</span></span>
        </div>
      </div>
    </div>

    <div class="section module-section">
      <div class="section-header">
        <h2>功能</h2>
        <span class="sub-title">同层级切换</span>
      </div>
      <div class="module-layer">
        <button class="module" :class="{ active: activeModule === 'daily' }" @click="activeModule = 'daily'">
          <h3>每日打卡</h3>
          <p>任务、积分、月活跃度、累计积分</p>
          <span class="num">01</span>
        </button>
        <button class="module" :class="{ active: activeModule === 'longTerm' }" @click="activeModule = 'longTerm'">
          <h3>长期主义</h3>
          <p>记录买了什么，计算实际日均消费</p>
          <span class="num">02</span>
        </button>
      </div>
    </div>

    <!-- 账户下方功能内容 -->
    <div v-if="activeModule === 'daily'" class="main-layout">
      <!-- 左侧：打卡任务 + 钱包 -->
      <div class="left-column">
        <!-- 打卡任务区域 -->
        <div class="section checkin-section">
          <div class="section-header">
            <h2>今日打卡任务</h2>
          </div>

          <div class="create-task">
            <form @submit.prevent="createTask">
              <input v-model="newTask.title" placeholder="任务标题" required />
              <select v-model="newTask.circle_mode">
                <option value="once">单次</option>
                <option value="weekly">每周一</option>
                <option value="workday">工作日</option>
                <option value="weekend">周末</option>
                <option value="custom">每天</option>
              </select>
              <select v-model.number="newTask.level">
                <option :value="1">● 低</option>
                <option :value="2">● 中</option>
                <option :value="3">● 高</option>
              </select>
              <input v-model.number="newTask.points" type="number" placeholder="积分" min="1" />
              <button type="submit" :disabled="creating">{{ creating ? '创建...' : '创建' }}</button>
            </form>
          </div>

          <div v-if="loadingTasks" class="loading">加载中...</div>
          <div v-else-if="todayTasks.length === 0" class="empty">今天没有需要打卡的任务</div>
          <div v-else class="tasks">
            <div class="task-item" v-for="task in paginatedTasks" :key="task.id" :class="['level-' + task.level, { 'checked': task.checked }]">
              <div class="task-info">
                <h4>
                  <span class="level-dot" :class="'dot-' + task.level"></span>
                  {{ task.title }}
                </h4>
                <p class="task-meta">
                  <span class="mode">{{ circleModeText(task.circle_mode) }}</span>
                  <span class="task-points">+{{ task.points }}</span>
                </p>
              </div>
              <div class="task-actions">
                <button class="edit-btn" @click="editTask(task)">修改</button>
                <button class="delete-btn" @click="deleteTask(task.id)">删除</button>
                <button v-if="!task.checked" class="skip-btn" @click="skipTask(task.id)">跳过</button>
                <button v-if="!task.checked" @click="checkin(task.id)">打卡</button>
                <button v-else class="cancel-btn" @click="cancelCheckin(task.id)">取消打卡</button>
              </div>
            </div>
            <!-- 分页按钮 -->
            <div v-if="totalPages > 1" class="pagination">
              <button class="page-btn" :disabled="currentPage === 1" @click="prevPage">上一页</button>
              <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
              <button class="page-btn" :disabled="currentPage === totalPages" @click="nextPage">下一页</button>
            </div>
          </div>
        </div>

        <!-- 积分钱包区域 -->
        <div class="section wallet-section">
          <div class="section-header">
            <h2>积分钱包</h2>
          </div>

          <div class="action-form">
            <form @submit.prevent="addRecord">
              <select v-model="addForm.type">
                <option value="earn">收入</option>
                <option value="spend">支出</option>
              </select>
              <input v-model.number="addForm.amount" type="number" placeholder="金额" min="1" required />
              <input v-model="addForm.description" placeholder="描述" required />
              <div class="date-wheel datetime-wheel">
                <select :value="getDateTimePart(addForm.record_time, 'year')" @change="setDateTimePart(addForm, 'record_time', 'year', $event.target.value)">
                  <option v-for="year in dateYears" :key="year" :value="year">{{ year }}年</option>
                </select>
                <select :value="getDateTimePart(addForm.record_time, 'month')" @change="setDateTimePart(addForm, 'record_time', 'month', $event.target.value)">
                  <option v-for="month in dateMonths" :key="month" :value="month">{{ month }}月</option>
                </select>
                <select :value="getDateTimePart(addForm.record_time, 'day')" @change="setDateTimePart(addForm, 'record_time', 'day', $event.target.value)">
                  <option v-for="day in dateDays(addForm.record_time)" :key="day" :value="day">{{ day }}日</option>
                </select>
                <select :value="getDateTimePart(addForm.record_time, 'hour')" @change="setDateTimePart(addForm, 'record_time', 'hour', $event.target.value)">
                  <option v-for="hour in timeHours" :key="hour" :value="hour">{{ pad2(hour) }}时</option>
                </select>
                <select :value="getDateTimePart(addForm.record_time, 'minute')" @change="setDateTimePart(addForm, 'record_time', 'minute', $event.target.value)">
                  <option v-for="minute in timeMinutes" :key="minute" :value="minute">{{ pad2(minute) }}分</option>
                </select>
              </div>
              <button type="submit" :disabled="adding">{{ adding ? '添加...' : '添加' }}</button>
            </form>
          </div>

          <div class="list-controls">
            <input v-model.trim="walletSearch" class="filter-input" placeholder="搜索钱包记录" />
          </div>

          <div v-if="loadingHistory" class="loading">加载中...</div>
          <div v-else-if="history.length === 0" class="empty">暂无记录</div>
          <div v-else-if="filteredHistory.length === 0" class="empty">没有匹配的钱包记录</div>
          <div v-else class="history-list">
            <div class="history-item" v-for="item in paginatedHistory" :key="item.id">
              <div class="item-info">
                <span :class="['type', item.type]">{{ item.type === 'earn' ? '收入' : '支出' }}</span>
                <span class="desc">{{ item.description }}</span>
                <span class="time">{{ formatDate(item.record_time) }}</span>
              </div>
              <div class="item-right">
                <span :class="['amount', item.type]">
                  {{ item.type === 'earn' ? '+' : '-' }}{{ item.amount }}
                </span>
                <button class="repeat-btn small" @click="repeatRecord(item)">重复</button>
                <button class="delete-btn small" @click="deleteRecord(item.id)">删除</button>
              </div>
            </div>
            <!-- 分页按钮 -->
            <div v-if="historyTotalPages > 1" class="pagination">
              <button class="page-btn" :disabled="historyPage === 1" @click="prevHistoryPage">上一页</button>
              <span class="page-info">{{ historyPage }} / {{ historyTotalPages }}</span>
              <button class="page-btn" :disabled="historyPage === historyTotalPages" @click="nextHistoryPage">下一页</button>
            </div>
          </div>
        </div>

        <!-- 其他任务 -->
        <div class="section other-tasks-section" v-if="otherTasks.length > 0">
          <div class="section-header">
            <h2>其他任务</h2>
            <span class="sub-title">今天不需要打卡</span>
          </div>
          <div class="tasks">
            <div class="task-item inactive" v-for="task in otherTasks" :key="task.id" :class="'level-' + task.level">
              <div class="task-info">
                <h4>
                  <span class="level-dot" :class="'dot-' + task.level"></span>
                  {{ task.title }}
                </h4>
                <p class="task-meta">
                  <span class="mode">{{ circleModeText(task.circle_mode) }}</span>
                  <span class="task-points">+{{ task.points }}</span>
                </p>
              </div>
              <div class="task-actions">
                <button class="edit-btn" @click="editTask(task)">修改</button>
                <button class="delete-btn" @click="deleteTask(task.id)">删除</button>
              </div>
            </div>
          </div>
        </div>
      </div>


      <!-- 编辑任务弹窗 -->
      <div v-if="showEditModal" class="modal-overlay" @click.self="closeEditModal">
        <div class="modal-content">
          <h3>修改任务</h3>
          <form @submit.prevent="updateTask">
            <div class="form-group">
              <label>任务标题</label>
              <input v-model="editForm.title" required />
            </div>
            <div class="form-group">
              <label>周期模式</label>
              <select v-model="editForm.circle_mode">
                <option value="once">单次</option>
                <option value="weekly">每周一</option>
                <option value="workday">工作日</option>
                <option value="weekend">周末</option>
                <option value="custom">每天</option>
              </select>
            </div>
            <div class="form-group">
              <label>任务级别</label>
              <select v-model.number="editForm.level">
                <option :value="1">● 低</option>
                <option :value="2">● 中</option>
                <option :value="3">● 高</option>
              </select>
            </div>
            <div class="form-group">
              <label>积分奖励</label>
              <input v-model.number="editForm.points" type="number" min="1" required />
            </div>
            <div class="modal-actions">
              <button type="button" class="cancel-btn" @click="closeEditModal">取消</button>
              <button type="submit" class="save-btn">保存</button>
            </div>
          </form>
        </div>
      </div>

      <!-- 右侧：热力图和统计图表 -->
      <div class="right-column">
        <div class="section heatmap-card">
          <div class="heatmap-header section-header">
            <h2>月活跃度</h2>
          </div>

          <div v-if="loadingHeatmap" class="loading">加载中...</div>
          <div v-else class="heatmap-board">
            <div class="heatmap-scroll">
              <div class="heatmap-content">
                <div class="heatmap-months" :style="{ gridTemplateColumns: `28px repeat(${heatmapWeeks.length}, minmax(0, 1fr))` }">
                  <span></span>
                  <span
                    v-for="(week, weekIndex) in heatmapWeeks"
                    :key="weekIndex"
                    class="heatmap-month-label"
                  >{{ week.showMonth ? week.monthLabel : '' }}</span>
                </div>
                <div class="heatmap-grid" :style="{ gridTemplateColumns: `28px repeat(${heatmapWeeks.length}, minmax(0, 1fr))` }">
                  <div class="heatmap-weekdays">
                    <span v-for="label in weekLabels" :key="label">{{ label }}</span>
                  </div>
                  <div class="heatmap-week-column" v-for="(week, weekIndex) in heatmapWeeks" :key="weekIndex" :title="`${week.month}月第${week.monthWeekNumber}周`">
                    <span
                      v-for="(day, dayIndex) in week.days"
                      :key="`${weekIndex}-${dayIndex}`"
                      class="heatmap-cell heatmap-day-cell"
                      :class="day ? heatmapClass(day.balance) : 'heatmap-empty'"
                      @mouseenter="showHeatmapTip(day)"
                      @mouseleave="clearHeatmapTip"
                      @click.stop="toggleHeatmapTip(day)"
                    >
                      <span v-if="day && activeHeatmapDay?.date === day.date" class="heatmap-tooltip">{{ formatHeatmapTitle(day) }}</span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
            <div class="heatmap-footer">
              <span>消耗</span>
              <span class="heatmap-cell neg-3"></span>
              <span class="heatmap-cell neg-2"></span>
              <span class="heatmap-cell neg-1"></span>
              <span class="heatmap-cell zero"></span>
              <span class="heatmap-cell pos-1"></span>
              <span class="heatmap-cell pos-2"></span>
              <span class="heatmap-cell pos-3"></span>
              <span>结余</span>
            </div>
          </div>
        </div>

        <div class="section chart-section">
          <div class="section-header">
            <h2>累计积分</h2>
            <select v-model.number="chartDays" @change="loadDailyStats" class="days-select">
              <option value="7">7天</option>
              <option value="14">14天</option>
              <option value="30">30天</option>
              <option value="180">180天</option>
              <option value="360">360天</option>
            </select>
          </div>

          <div class="chart-wrapper">
            <div v-if="loadingStats" class="loading">加载中...</div>
            <canvas v-show="!loadingStats" ref="chartCanvas"></canvas>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="long-term-layout">
      <div class="section long-term-section">
        <div class="section-header">
          <h2>长期主义</h2>
          <span class="sub-title">当前日均消费 {{ formatDayCost(longTermSummary.active_daily_cost) }}</span>
        </div>

        <div class="long-term-summary">
          <div class="stat-card">
            <label>日均消费</label>
            <strong>{{ formatDayCost(longTermSummary.active_daily_cost) }}</strong>
          </div>
          <div class="stat-card">
            <label>使用中</label>
            <strong>{{ longTermSummary.active_count }} 件</strong>
          </div>
          <div class="stat-card">
            <label>已报废</label>
            <strong>{{ longTermSummary.scrapped_count }} 件</strong>
          </div>
        </div>

        <div class="action-form long-term-form">
          <form @submit.prevent="createLongTermItem">
            <input v-model.trim="longTermForm.name" placeholder="买了什么" required />
            <input v-model.number="longTermForm.price" type="number" placeholder="价格" min="0.01" step="0.01" required />
            <div class="date-wheel">
              <select :value="getDatePart(longTermForm.purchase_date, 'year')" @change="setDatePart(longTermForm, 'purchase_date', 'year', $event.target.value)">
                <option v-for="year in dateYears" :key="year" :value="year">{{ year }}年</option>
              </select>
              <select :value="getDatePart(longTermForm.purchase_date, 'month')" @change="setDatePart(longTermForm, 'purchase_date', 'month', $event.target.value)">
                <option v-for="month in dateMonths" :key="month" :value="month">{{ month }}月</option>
              </select>
              <select :value="getDatePart(longTermForm.purchase_date, 'day')" @change="setDatePart(longTermForm, 'purchase_date', 'day', $event.target.value)">
                <option v-for="day in dateDays(longTermForm.purchase_date)" :key="day" :value="day">{{ day }}日</option>
              </select>
            </div>
            <button type="submit" :disabled="creatingLongTerm">{{ creatingLongTerm ? '记录中...' : '记录' }}</button>
          </form>
        </div>

        <div v-if="loadingLongTerm" class="loading">加载中...</div>
        <template v-else>
          <div class="long-term-list-block">
            <div class="section-header inline-header">
              <h2>使用中</h2>
              <span class="sub-title">日均成本会随使用天数下降</span>
            </div>
            <div v-if="allActiveLongTermItems.length > 0" class="list-controls long-term-controls">
              <input v-model.trim="longTermSearch" class="filter-input" placeholder="搜索使用中物品" />
              <select v-model="longTermSort" class="sort-select">
                <option value="date_asc">时间顺序</option>
                <option value="date_desc">时间倒序</option>
                <option value="price_asc">价格顺序</option>
                <option value="price_desc">价格倒序</option>
                <option value="daily_cost_desc">日均价格倒序</option>
              </select>
            </div>
            <div v-if="allActiveLongTermItems.length === 0" class="empty">还没有使用中的物品</div>
            <div v-else-if="activeLongTermItems.length === 0" class="empty">没有匹配的使用中物品</div>
            <div v-else class="item-list">
              <div class="long-term-item" v-for="item in activeLongTermItems" :key="item.id">
                <template v-if="editingLongTermId === item.id">
                  <div class="long-term-edit-form">
                    <input v-model.trim="editLongTermForm.name" placeholder="买了什么" required />
                    <input v-model.number="editLongTermForm.price" type="number" placeholder="价格" min="0.01" step="0.01" required />
                    <div class="date-wheel">
                      <select :value="getDatePart(editLongTermForm.purchase_date, 'year')" @change="setDatePart(editLongTermForm, 'purchase_date', 'year', $event.target.value)">
                        <option v-for="year in dateYears" :key="year" :value="year">{{ year }}年</option>
                      </select>
                      <select :value="getDatePart(editLongTermForm.purchase_date, 'month')" @change="setDatePart(editLongTermForm, 'purchase_date', 'month', $event.target.value)">
                        <option v-for="month in dateMonths" :key="month" :value="month">{{ month }}月</option>
                      </select>
                      <select :value="getDatePart(editLongTermForm.purchase_date, 'day')" @change="setDatePart(editLongTermForm, 'purchase_date', 'day', $event.target.value)">
                        <option v-for="day in dateDays(editLongTermForm.purchase_date)" :key="day" :value="day">{{ day }}日</option>
                      </select>
                    </div>
                    <div class="long-term-actions">
                      <button class="edit-btn" @click="saveLongTermItem(item)">保存</button>
                      <button class="skip-btn" @click="cancelEditLongTermItem">取消</button>
                    </div>
                  </div>
                </template>
                <template v-else>
                  <div class="item-icon">{{ item.name.slice(0, 1) }}</div>
                  <div class="long-term-main">
                    <h4>{{ item.name }}</h4>
                    <p>{{ formatMoney(item.price) }} · {{ formatDayCost(item.daily_cost) }}</p>
                  </div>
                  <div class="long-term-side">
                    <span>{{ item.owned_days }}天</span>
                    <div class="long-term-actions">
                      <button class="edit-btn" @click="editLongTermItem(item)">修改</button>
                      <button class="skip-btn" @click="scrapLongTermItem(item)">报废</button>
                      <button class="delete-btn" @click="deleteLongTermItem(item)">删除</button>
                    </div>
                  </div>
                </template>
              </div>
            </div>
          </div>

          <div class="long-term-list-block scrapped-block">
            <div class="section-header inline-header">
              <h2>报废栏</h2>
              <span class="sub-title">日均成本已冻结</span>
            </div>
            <div v-if="scrappedLongTermItems.length === 0" class="empty">暂无报废物品</div>
            <div v-else class="item-list">
              <div class="long-term-item scrapped" v-for="item in scrappedLongTermItems" :key="item.id">
                <template v-if="editingLongTermId === item.id">
                  <div class="long-term-edit-form">
                    <input v-model.trim="editLongTermForm.name" placeholder="买了什么" required />
                    <input v-model.number="editLongTermForm.price" type="number" placeholder="价格" min="0.01" step="0.01" required />
                    <div class="date-wheel">
                      <select :value="getDatePart(editLongTermForm.purchase_date, 'year')" @change="setDatePart(editLongTermForm, 'purchase_date', 'year', $event.target.value)">
                        <option v-for="year in dateYears" :key="year" :value="year">{{ year }}年</option>
                      </select>
                      <select :value="getDatePart(editLongTermForm.purchase_date, 'month')" @change="setDatePart(editLongTermForm, 'purchase_date', 'month', $event.target.value)">
                        <option v-for="month in dateMonths" :key="month" :value="month">{{ month }}月</option>
                      </select>
                      <select :value="getDatePart(editLongTermForm.purchase_date, 'day')" @change="setDatePart(editLongTermForm, 'purchase_date', 'day', $event.target.value)">
                        <option v-for="day in dateDays(editLongTermForm.purchase_date)" :key="day" :value="day">{{ day }}日</option>
                      </select>
                    </div>
                    <div class="long-term-actions">
                      <button class="edit-btn" @click="saveLongTermItem(item)">保存</button>
                      <button class="skip-btn" @click="cancelEditLongTermItem">取消</button>
                    </div>
                  </div>
                </template>
                <template v-else>
                  <div class="item-icon">{{ item.name.slice(0, 1) }}</div>
                  <div class="long-term-main">
                    <h4>{{ item.name }}</h4>
                    <p>{{ formatMoney(item.price) }} · {{ formatDayCost(item.daily_cost) }}</p>
                  </div>
                  <div class="long-term-side">
                    <span>{{ item.owned_days }}天</span>
                    <div class="long-term-actions">
                      <button class="edit-btn" @click="editLongTermItem(item)">修改</button>
                      <button class="delete-btn" @click="deleteLongTermItem(item)">删除</button>
                    </div>
                  </div>
                </template>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 编辑资料弹窗 -->
    <div v-if="showProfileModal" class="modal-overlay" @click.self="closeProfileModal">
      <div class="modal-content profile-modal">
        <h3>编辑资料</h3>
        <form @submit.prevent="saveProfile">
          <div class="avatar-edit">
            <div class="profile-avatar large">
              <img v-if="avatarPreview || profileForm.avatar_url" :src="avatarPreview || profileForm.avatar_url" alt="头像预览" />
              <span v-else>{{ profileInitial }}</span>
            </div>
            <label class="avatar-upload-btn">
              选择头像
              <input type="file" accept="image/jpeg,image/png,image/webp,image/gif" @change="handleAvatarChange" />
            </label>
          </div>
          <div class="form-group">
            <label>用户名</label>
            <input v-model.trim="profileForm.username" type="text" required />
          </div>
          <div class="form-group">
            <label>邮箱</label>
            <input v-model.trim="profileForm.email" type="email" placeholder="可选" />
          </div>
          <div class="form-group">
            <label>旧密码</label>
            <input v-model="profileForm.old_password" type="password" placeholder="修改密码时必填" />
          </div>
          <div class="form-group">
            <label>新密码</label>
            <input v-model="profileForm.new_password" type="password" placeholder="不修改可留空" />
          </div>
          <div class="form-group">
            <label>确认新密码</label>
            <input v-model="profileForm.confirm_password" type="password" placeholder="再次输入新密码" />
          </div>
          <p class="error" v-if="profileError">{{ profileError }}</p>
          <p class="success" v-if="profileSuccess">{{ profileSuccess }}</p>
          <div class="modal-actions">
            <button type="button" class="cancel-btn" @click="closeProfileModal">取消</button>
            <button type="submit" class="save-btn" :disabled="savingProfile">{{ savingProfile ? '保存中...' : '保存' }}</button>
          </div>
        </form>
      </div>
    </div>

    <footer class="footer">
      <span>© 2026 srDailyTask</span>
      <a href="https://github.com/MGter/srDailyTask" target="_blank">GitHub</a>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { taskApi, checkinApi, userApi, walletApi, pointsApi, longTermApi } from '../api'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

const router = useRouter()
const user = ref(null)
const tasks = ref([])
const history = ref([])
const activeModule = ref('daily')
const longTermItems = ref([])
const longTermSummary = ref({ active_daily_cost: 0, active_count: 0, scrapped_count: 0 })
const loadingLongTerm = ref(false)
const creatingLongTerm = ref(false)
const editingLongTermId = ref(null)
const editLongTermForm = ref({ name: '', price: null, purchase_date: '' })
const walletSearch = ref('')
const longTermSearch = ref('')
const longTermSort = ref('daily_cost_desc')
const historyPage = ref(1)
const historyPageSize = 5

const normalizeSearchText = (value) => String(value ?? '').toLowerCase()

const walletRecordSearchText = (item) => [
  item.type,
  item.type === 'earn' ? '收入' : '支出',
  item.description,
  item.amount,
  item.record_time,
  formatDate(item.record_time)
].join(' ')

const filteredHistory = computed(() => {
  const keyword = normalizeSearchText(walletSearch.value.trim())
  if (!keyword) return history.value
  return history.value.filter(item => normalizeSearchText(walletRecordSearchText(item)).includes(keyword))
})

const historyTotalPages = computed(() => Math.ceil(filteredHistory.value.length / historyPageSize))

const paginatedHistory = computed(() => {
  const start = (historyPage.value - 1) * historyPageSize
  return filteredHistory.value.slice(start, start + historyPageSize)
})

watch(walletSearch, () => {
  historyPage.value = 1
})

watch(filteredHistory, () => {
  if (historyTotalPages.value > 0 && historyPage.value > historyTotalPages.value) {
    historyPage.value = historyTotalPages.value
  }
})

const nextHistoryPage = () => {
  if (historyPage.value < historyTotalPages.value) historyPage.value++
}

const prevHistoryPage = () => {
  if (historyPage.value > 1) historyPage.value--
}
const dailyStats = ref([])
const heatmapStats = ref([])
const loadingTasks = ref(true)
const loadingHistory = ref(false)
const loadingStats = ref(false)
const loadingHeatmap = ref(false)
const creating = ref(false)
const adding = ref(false)
const chartDays = ref(7)
const chartCanvas = ref(null)
let chartInstance = null

const userId = parseInt(localStorage.getItem('userId'))

const newTask = ref({
  title: '',
  circle_mode: 'workday',
  level: 1,
  points: 10
})

const showEditModal = ref(false)
const editForm = ref({
  id: 0,
  title: '',
  circle_mode: 'workday',
  level: 1,
  points: 10
})

const getDefaultTime = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

const addForm = ref({
  type: 'spend',
  amount: 10,
  description: '',
  record_time: getDefaultTime()
})

const localDateInput = (value = new Date()) => {
  if (typeof value === 'string' && /^\d{4}-\d{2}-\d{2}/.test(value)) {
    return value.slice(0, 10)
  }
  const date = value instanceof Date ? value : new Date(value)
  const safeDate = Number.isNaN(date.getTime()) ? new Date() : date
  return `${safeDate.getFullYear()}-${String(safeDate.getMonth() + 1).padStart(2, '0')}-${String(safeDate.getDate()).padStart(2, '0')}`
}

const longTermForm = ref({
  name: '',
  price: null,
  purchase_date: localDateInput()
})

const dateYears = computed(() => {
  const currentYear = new Date().getFullYear()
  const years = []
  for (let year = currentYear + 1; year >= 2000; year--) years.push(year)
  return years
})
const dateMonths = Array.from({ length: 12 }, (_, index) => index + 1)
const timeHours = Array.from({ length: 24 }, (_, index) => index)
const timeMinutes = Array.from({ length: 60 }, (_, index) => index)

const pad2 = (value) => String(value).padStart(2, '0')
const daysInMonth = (year, month) => new Date(Number(year), Number(month), 0).getDate()
const formValue = (target) => target?.value || target
const getDatePart = (value, part) => {
  const [year, month, day] = localDateInput(value).split('-').map(Number)
  if (part === 'year') return year
  if (part === 'month') return month
  return day
}
const getDateTimePart = (value, part) => {
  const current = String(value || getDefaultTime())
  if (part === 'hour') return Number(current.slice(11, 13))
  if (part === 'minute') return Number(current.slice(14, 16))
  return getDatePart(current, part)
}
const dateDays = (value) => Array.from({ length: daysInMonth(getDatePart(value, 'year'), getDatePart(value, 'month')) }, (_, index) => index + 1)
const setDatePart = (target, field, part, rawValue) => {
  const form = formValue(target)
  const value = Number(rawValue)
  const year = part === 'year' ? value : getDatePart(form[field], 'year')
  const month = part === 'month' ? value : getDatePart(form[field], 'month')
  const maxDay = daysInMonth(year, month)
  const day = Math.min(part === 'day' ? value : getDatePart(form[field], 'day'), maxDay)
  form[field] = `${year}-${pad2(month)}-${pad2(day)}`
}
const setDateTimePart = (target, field, part, rawValue) => {
  const form = formValue(target)
  const current = form[field] || getDefaultTime()
  const value = Number(rawValue)
  const year = part === 'year' ? value : getDateTimePart(current, 'year')
  const month = part === 'month' ? value : getDateTimePart(current, 'month')
  const maxDay = daysInMonth(year, month)
  const day = Math.min(part === 'day' ? value : getDateTimePart(current, 'day'), maxDay)
  const hour = part === 'hour' ? value : getDateTimePart(current, 'hour')
  const minute = part === 'minute' ? value : getDateTimePart(current, 'minute')
  form[field] = `${year}-${pad2(month)}-${pad2(day)}T${pad2(hour)}:${pad2(minute)}`
}
const apiDate = (value) => new Date(`${localDateInput(value)}T00:00:00`).toISOString()
const itemDateInput = (value) => localDateInput(value)

const showProfileModal = ref(false)
const savingProfile = ref(false)
const profileError = ref('')
const profileSuccess = ref('')
const avatarFile = ref(null)
const avatarPreview = ref('')
const profileForm = ref({
  username: '',
  email: '',
  avatar_url: '',
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const todayText = computed(() => {
  const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const now = new Date()
  return `${now.getFullYear()}年${now.getMonth()+1}月${now.getDate()}日 ${days[now.getDay()]}`
})

const todayStats = computed(() => {
  const today = localDateInput()
  const stat = dailyStats.value.find(s => s.date === today)
  return {
    earn: stat?.earn || 0,
    spend: stat?.spend || 0,
    balance: stat?.balance || 0,
    total: user.value?.points || 0
  }
})

const userInitial = computed(() => user.value?.username?.slice(0, 1).toUpperCase() || '?')
const profileInitial = computed(() => profileForm.value.username?.slice(0, 1).toUpperCase() || userInitial.value)
const allActiveLongTermItems = computed(() => longTermItems.value.filter(item => item.status === 'active'))
const scrappedLongTermItems = computed(() => longTermItems.value.filter(item => item.status === 'scrapped'))
const formatMoney = (value) => `¥${Number(value || 0).toFixed(2)}`
const formatDayCost = (value) => `${formatMoney(value)}/天`

const longTermItemSearchText = (item) => [
  item.name,
  item.price,
  formatMoney(item.price),
  item.purchase_date,
  localDateInput(item.purchase_date),
  item.daily_cost,
  formatDayCost(item.daily_cost),
  item.owned_days,
  `${item.owned_days}天`,
  item.status,
  item.status === 'active' ? '使用中' : '报废'
].join(' ')

const activeLongTermItems = computed(() => {
  const keyword = normalizeSearchText(longTermSearch.value.trim())
  const filtered = keyword
    ? allActiveLongTermItems.value.filter(item => normalizeSearchText(longTermItemSearchText(item)).includes(keyword))
    : allActiveLongTermItems.value

  return [...filtered].sort((a, b) => {
    if (longTermSort.value === 'date_asc') return new Date(a.purchase_date) - new Date(b.purchase_date)
    if (longTermSort.value === 'date_desc') return new Date(b.purchase_date) - new Date(a.purchase_date)
    if (longTermSort.value === 'price_asc') return Number(a.price) - Number(b.price)
    if (longTermSort.value === 'price_desc') return Number(b.price) - Number(a.price)
    return Number(b.daily_cost) - Number(a.daily_cost)
  })
})

const weekLabels = ['一', '二', '三', '四', '五', '六', '日']
const monthLabels = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
const activeHeatmapDay = ref(null)

const formatHeatmapTitle = (day) => {
  if (!day) return ''
  const weekText = day.monthWeekNumber ? `${day.month}月第${day.monthWeekNumber}周\n` : ''
  return `${weekText}${day.date}\n累计 +${day.earn}\n消耗 -${day.spend}\n结余 ${day.balance >= 0 ? '+' : ''}${day.balance}`
}

const showHeatmapTip = (day) => {
  if (day) activeHeatmapDay.value = day
}

const clearHeatmapTip = () => {
  activeHeatmapDay.value = null
}

const toggleHeatmapTip = (day) => {
  if (!day) return
  activeHeatmapDay.value = activeHeatmapDay.value?.date === day.date ? null : day
}

const formatHeatmapDate = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const heatmapPeriodDays = computed(() => {
  const statsByDate = new Map(heatmapStats.value.map(s => [s.date, s]))
  const end = new Date()
  end.setHours(0, 0, 0, 0)
  const start = new Date(end)
  start.setMonth(start.getMonth() - 2)
  start.setDate(start.getDate() + 1)

  const days = []
  for (const day = new Date(start); day <= end; day.setDate(day.getDate() + 1)) {
    const date = formatHeatmapDate(day)
    const stat = statsByDate.get(date)
    days.push({
      date,
      earn: stat?.earn || 0,
      spend: stat?.spend || 0,
      balance: stat?.balance || 0
    })
  }
  return days
})

const getMonthWeekNumber = (date) => {
  const firstDay = new Date(date.getFullYear(), date.getMonth(), 1)
  const firstWeekday = (firstDay.getDay() + 6) % 7
  return Math.floor((date.getDate() + firstWeekday - 1) / 7) + 1
}

const heatmapClass = (balance) => {
  if (balance >= 100) return 'pos-3'
  if (balance >= 50) return 'pos-2'
  if (balance > 0) return 'pos-1'
  if (balance <= -100) return 'neg-3'
  if (balance <= -50) return 'neg-2'
  if (balance < 0) return 'neg-1'
  return 'zero'
}

const heatmapWeeks = computed(() => {
  const stats = heatmapPeriodDays.value.map(s => {
    const date = new Date(`${s.date}T00:00:00`)
    return {
      date: s.date,
      earn: s.earn,
      spend: s.spend,
      balance: s.balance,
      weekday: (date.getDay() + 6) % 7,
      month: Number(s.date.slice(5, 7)),
      monthWeekNumber: getMonthWeekNumber(date)
    }
  })

  if (stats.length === 0) return []

  const weeks = []
  let currentWeek = Array(7).fill(null)

  stats.forEach((day, index) => {
    currentWeek[day.weekday] = day

    const isLastDay = index === stats.length - 1
    if (day.weekday === 6 || isLastDay) {
      const firstDay = currentWeek.find(Boolean)
      const previousWeek = weeks[weeks.length - 1]
      const weekNumber = weeks.length + 1
      weeks.push({
        days: currentWeek.map(item => item ? { ...item, weekNumber } : null),
        month: firstDay?.month || day.month,
        monthWeekNumber: firstDay?.monthWeekNumber || day.monthWeekNumber,
        monthLabel: monthLabels[(firstDay?.month || day.month) - 1],
        showMonth: !previousWeek || previousWeek.month !== (firstDay?.month || day.month),
        weekNumber
      })
      currentWeek = Array(7).fill(null)
    }
  })

  return weeks
})

const todayTasks = computed(() => {
  const filtered = tasks.value.filter(t => t.should_checkin_today && !t.is_expired)
  // 未打卡任务排在前面，已打卡排在后面
  return filtered.sort((a, b) => {
    if (a.checked === b.checked) return 0
    return a.checked ? 1 : -1
  })
})

const pageSize = 5
const currentPage = ref(1)

const totalPages = computed(() => Math.ceil(todayTasks.value.length / pageSize))

const paginatedTasks = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return todayTasks.value.slice(start, start + pageSize)
})

const nextPage = () => {
  if (currentPage.value < totalPages.value) currentPage.value++
}

const prevPage = () => {
  if (currentPage.value > 1) currentPage.value--
}

const otherTasks = computed(() => {
  return tasks.value.filter(t => !t.should_checkin_today && !t.is_expired)
})

const circleModeText = (mode) => {
  const texts = {
    once: '单次',
    weekly: '每周一',
    workday: '工作日',
    weekend: '周末',
    custom: '每天'
  }
  return texts[mode] || mode
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const logout = () => {
  localStorage.removeItem('userId')
  localStorage.removeItem('user')
  router.push('/login')
}

const clearAvatarPreview = () => {
  if (avatarPreview.value) {
    URL.revokeObjectURL(avatarPreview.value)
    avatarPreview.value = ''
  }
}

const openProfileModal = () => {
  profileForm.value = {
    username: user.value?.username || '',
    email: user.value?.email || '',
    avatar_url: user.value?.avatar_url || '',
    old_password: '',
    new_password: '',
    confirm_password: ''
  }
  avatarFile.value = null
  clearAvatarPreview()
  profileError.value = ''
  profileSuccess.value = ''
  showProfileModal.value = true
}

const closeProfileModal = () => {
  showProfileModal.value = false
  avatarFile.value = null
  clearAvatarPreview()
}

const handleAvatarChange = (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    profileError.value = '头像只支持 jpg、png、webp、gif'
    return
  }
  if (file.size > 2 * 1024 * 1024) {
    profileError.value = '头像不能超过 2MB'
    return
  }

  avatarFile.value = file
  clearAvatarPreview()
  avatarPreview.value = URL.createObjectURL(file)
  profileError.value = ''
}

const saveProfile = async () => {
  profileError.value = ''
  profileSuccess.value = ''
  if (profileForm.value.new_password !== profileForm.value.confirm_password) {
    profileError.value = '两次输入的新密码不一致'
    return
  }

  savingProfile.value = true
  try {
    const payload = {
      username: profileForm.value.username,
      email: profileForm.value.email,
      old_password: profileForm.value.old_password,
      new_password: profileForm.value.new_password
    }
    let res = await userApi.updateUser(userId, payload)

    if (avatarFile.value) {
      const formData = new FormData()
      formData.append('avatar', avatarFile.value)
      res = await userApi.uploadAvatar(userId, formData)
    }

    user.value = res.data
    localStorage.setItem('user', JSON.stringify(res.data))
    profileSuccess.value = '资料已更新'
    setTimeout(() => {
      closeProfileModal()
    }, 600)
  } catch (e) {
    profileError.value = e.response?.data?.error || '保存失败'
  } finally {
    savingProfile.value = false
  }
}

const loadUser = async () => {
  if (!userId) {
    router.push('/login')
    return
  }
  try {
    const res = await userApi.getUser(userId)
    user.value = res.data
    localStorage.setItem('user', JSON.stringify(res.data))
  } catch (e) {
    router.push('/login')
  }
}

const loadTasks = async () => {
  loadingTasks.value = true
  try {
    const res = await taskApi.getUserTasks(userId, { limit: 50, offset: 0 })
    const checkedRes = await checkinApi.getTodayChecked(userId)
    const checkedIds = checkedRes.data || []
    tasks.value = res.data.map(task => ({
      ...task,
      checked: checkedIds.includes(task.id)
    }))
  } catch (e) {
    console.error('加载任务失败', e)
  } finally {
    loadingTasks.value = false
  }
}

const createTask = async () => {
  creating.value = true
  try {
    const res = await taskApi.create({
      user_id: userId,
      title: newTask.value.title,
      circle_mode: newTask.value.circle_mode,
      level: Number(newTask.value.level),
      points: Number(newTask.value.points) || 10
    })
    const checkedRes = await checkinApi.getTodayChecked(userId)
    const checkedIds = checkedRes.data || []
    tasks.value.unshift({
      ...res.data,
      checked: checkedIds.includes(res.data.id)
    })
    newTask.value.title = ''
  } catch (e) {
    alert('创建失败：' + (e.response?.data?.error || '未知错误'))
  } finally {
    creating.value = false
  }
}

const deleteTask = async (taskId) => {
  if (!confirm('确定删除该任务吗？')) return
  try {
    await taskApi.delete(taskId)
    tasks.value = tasks.value.filter(t => t.id !== taskId)
  } catch (e) {
    alert('删除失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const editTask = (task) => {
  editForm.value = {
    id: task.id,
    title: task.title,
    circle_mode: task.circle_mode,
    level: task.level,
    points: task.points
  }
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
}

const updateTask = async () => {
  try {
    await taskApi.update(editForm.value.id, {
      title: editForm.value.title,
      circle_mode: editForm.value.circle_mode,
      level: Number(editForm.value.level),
      points: Number(editForm.value.points)
    })
    // 更新本地任务列表
    const task = tasks.value.find(t => t.id === editForm.value.id)
    if (task) {
      task.title = editForm.value.title
      task.circle_mode = editForm.value.circle_mode
      task.level = editForm.value.level
      task.points = editForm.value.points
    }
    closeEditModal()
  } catch (e) {
    alert('修改失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const checkin = async (taskId) => {
  try {
    await checkinApi.checkin(taskId, { user_id: userId })
    const task = tasks.value.find(t => t.id === taskId)
    if (task) task.checked = true
    if (task && task.circle_mode === 'once') {
      task.is_expired = true
    }
    const res = await userApi.getUser(userId)
    user.value = res.data
    await loadHistory()
    await refreshStats()
  } catch (e) {
    alert('打卡失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const cancelCheckin = async (taskId) => {
  if (!confirm('确定取消打卡吗？积分将退还。')) return
  try {
    await checkinApi.cancel(taskId, { user_id: userId })
    const task = tasks.value.find(t => t.id === taskId)
    if (task) task.checked = false
    if (task && task.circle_mode === 'once') {
      task.is_expired = false
    }
    const res = await userApi.getUser(userId)
    user.value = res.data
    await loadHistory()
    await refreshStats()
  } catch (e) {
    alert('取消失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const skipTask = async (taskId) => {
  try {
    await checkinApi.skip(taskId, { user_id: userId })
    const task = tasks.value.find(t => t.id === taskId)
    if (task) task.checked = true
    await loadHistory()
    await refreshStats()
  } catch (e) {
    alert('跳过失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const loadHistory = async () => {
  loadingHistory.value = true
  try {
    const res = await pointsApi.getHistory(userId, { limit: 30, offset: 0 })
    history.value = res.data
  } catch (e) {
    console.error('加载记录失败', e)
  } finally {
    loadingHistory.value = false
  }
}

const loadLongTermItems = async () => {
  loadingLongTerm.value = true
  try {
    const res = await longTermApi.list(userId)
    longTermItems.value = res.data.items || []
    longTermSummary.value = res.data.summary || { active_daily_cost: 0, active_count: 0, scrapped_count: 0 }
  } catch (e) {
    console.error('加载长期主义失败', e)
  } finally {
    loadingLongTerm.value = false
  }
}

const createLongTermItem = async () => {
  creatingLongTerm.value = true
  try {
    await longTermApi.create({
      user_id: userId,
      name: longTermForm.value.name,
      price: Number(longTermForm.value.price),
      purchase_date: apiDate(longTermForm.value.purchase_date)
    })
    longTermForm.value = {
      name: '',
      price: null,
      purchase_date: localDateInput()
    }
    await loadLongTermItems()
  } catch (e) {
    alert('记录失败：' + (e.response?.data?.error || '未知错误'))
  } finally {
    creatingLongTerm.value = false
  }
}

const editLongTermItem = (item) => {
  editingLongTermId.value = item.id
  editLongTermForm.value = {
    name: item.name,
    price: item.price,
    purchase_date: itemDateInput(item.purchase_date)
  }
}

const cancelEditLongTermItem = () => {
  editingLongTermId.value = null
}

const saveLongTermItem = async (item) => {
  try {
    await longTermApi.update(item.id, {
      user_id: userId,
      name: editLongTermForm.value.name,
      price: Number(editLongTermForm.value.price),
      purchase_date: apiDate(editLongTermForm.value.purchase_date)
    })
    editingLongTermId.value = null
    await loadLongTermItems()
  } catch (e) {
    alert('修改失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const scrapLongTermItem = async (item) => {
  if (!confirm(`确定报废「${item.name}」吗？日均成本会冻结。`)) return
  try {
    await longTermApi.scrap(item.id, {
      user_id: userId,
      scrap_date: new Date().toISOString()
    })
    await loadLongTermItems()
  } catch (e) {
    alert('报废失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const deleteLongTermItem = async (item) => {
  if (!confirm(`确定删除「${item.name}」吗？`)) return
  try {
    await longTermApi.delete(item.id, userId)
    await loadLongTermItems()
  } catch (e) {
    alert('删除失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const addRecord = async () => {
  adding.value = true
  try {
    await walletApi.addRecord({
      user_id: userId,
      type: addForm.value.type,
      amount: addForm.value.amount,
      description: addForm.value.description,
      record_time: new Date(addForm.value.record_time).toISOString()
    })
    const res = await userApi.getUser(userId)
    user.value = res.data
    await loadHistory()
    await refreshStats()
    addForm.value.description = ''
    addForm.value.record_time = getDefaultTime()
  } catch (e) {
    alert('添加失败：' + (e.response?.data?.error || '未知错误'))
  } finally {
    adding.value = false
  }
}

const deleteRecord = async (id) => {
  if (!confirm('确定删除该记录吗？删除后会相应调整积分余额。')) return
  try {
    await walletApi.delete(id, userId)
    const res = await userApi.getUser(userId)
    user.value = res.data
    await loadHistory()
    await refreshStats()
  } catch (e) {
    alert('删除失败：' + (e.response?.data?.error || '未知错误'))
  }
}

const repeatRecord = (item) => {
  addForm.value.type = item.type
  addForm.value.amount = item.amount
  addForm.value.description = item.description
  addForm.value.record_time = getDefaultTime()  // 使用当前时间
}

const loadDailyStats = async () => {
  loadingStats.value = true
  try {
    const res = await pointsApi.getDailyStats(userId, { days: chartDays.value })
    dailyStats.value = res.data
    // 等待DOM更新后再渲染图表
    await nextTick()
    // 再次等待确保canvas已渲染
    setTimeout(() => {
      renderChart()
    }, 100)
  } catch (e) {
    console.error('加载统计失败', e)
  } finally {
    loadingStats.value = false
  }
}

const loadHeatmapStats = async () => {
  loadingHeatmap.value = true
  try {
    const res = await pointsApi.getDailyStats(userId, { days: 62 })
    heatmapStats.value = res.data
  } catch (e) {
    console.error('加载热力图失败', e)
  } finally {
    loadingHeatmap.value = false
  }
}

const refreshStats = async () => {
  await loadDailyStats()
  await loadHeatmapStats()
}

const renderChart = () => {
  console.log('renderChart called, canvas:', chartCanvas.value, 'data:', dailyStats.value.length)
  if (!chartCanvas.value) {
    console.log('No canvas found')
    return
  }
  if (dailyStats.value.length === 0) {
    console.log('No data')
    return
  }

  if (chartInstance) {
    chartInstance.destroy()
  }

  const ctx = chartCanvas.value.getContext('2d')
  const labels = dailyStats.value.map(s => s.date.slice(5))

  // 累计积分：从第一天开始累加每日balance
  let cumulative = 0
  const cumulativeData = dailyStats.value.map(s => {
    cumulative += s.balance
    return cumulative
  })

  const cumulativeFillColor = (context) => {
    const chart = context.chart
    const { ctx, chartArea, scales } = chart
    if (!chartArea) return 'rgba(52, 199, 89, 0.18)'

    const zeroY = scales.y.getPixelForValue(0)
    const offset = Math.max(0, Math.min(1, (zeroY - chartArea.top) / (chartArea.bottom - chartArea.top)))
    const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom)
    gradient.addColorStop(0, 'rgba(52, 199, 89, 0.42)')
    gradient.addColorStop(Math.max(0, offset - 0.16), 'rgba(52, 199, 89, 0.24)')
    gradient.addColorStop(offset, 'rgba(52, 199, 89, 0.04)')
    gradient.addColorStop(offset, 'rgba(255, 59, 48, 0.04)')
    gradient.addColorStop(Math.min(1, offset + 0.16), 'rgba(255, 59, 48, 0.24)')
    gradient.addColorStop(1, 'rgba(255, 59, 48, 0.42)')
    return gradient
  }

  console.log('Chart labels:', labels, 'cumulativeData:', cumulativeData)

  chartInstance = new Chart(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [
        {
          label: '累计积分',
          data: cumulativeData,
          borderColor: 'rgba(0, 0, 0, 0)',
          borderWidth: 0,
          backgroundColor: cumulativeFillColor,
          tension: 0.3,
          pointRadius: 0,
          pointHoverRadius: 0,
          fill: 'origin'
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        }
      },
      scales: {
        y: {
          position: 'left',
          beginAtZero: true,
          ticks: {
            color: '#1d1d1f'
          },
          grid: {
            color: '#e5e5e5'
          }
        },
        x: {
          ticks: {
            color: '#1d1d1f'
          },
          grid: {
            display: false
          }
        }
      }
    }
  })
}

watch(activeModule, async (module) => {
  if (module === 'daily') {
    await nextTick()
    setTimeout(() => {
      renderChart()
    }, 50)
    return
  }
  if (chartInstance) {
    chartInstance.destroy()
    chartInstance = null
  }
})

onMounted(async () => {
  await loadUser()
  await loadLongTermItems()
  await loadTasks()
  await loadHistory()
  await refreshStats()
})

onUnmounted(() => {
  clearAvatarPreview()
})
</script>

<style scoped>
* {
  box-sizing: border-box;
}

.dashboard-container {
  position: relative;
  z-index: 0;
  width: min(100%, 1120px);
  max-width: 100%;
  min-height: 100vh;
  margin: 0 auto;
  padding: 20px 15px 40px 15px;
  background: transparent;
}

.dashboard-container::before {
  content: '';
  position: fixed;
  inset: 0;
  z-index: -1;
  background-image: url('/assets/kita.webp');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  pointer-events: none;
}

.dashboard-container > * {
  position: relative;
  z-index: 1;
}

/* 手机端调整 */

.points-cards {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 8px;
  align-items: stretch;
}

.point-card {
  min-width: 0;
  padding: 10px 8px;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.52);
  color: #1d1d1f;
  font-size: 12px;
  font-weight: 700;
  text-align: center;
  white-space: nowrap;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.52);
}

.point-label,
.point-value {
  display: inline;
}

.point-card.earn { color: #1f8f45; }
.point-card.spend { color: #c92a22; }
.point-card.balance { color: #b26a00; }
.point-card.total { color: #0066d6; }
.point-card.daily-cost { color: #7c3aed; }

.logout-btn {
  padding: 7px 10px;
  background: rgba(255, 59, 48, 0.12);
  color: #d70015;
  border: 1px solid rgba(255, 59, 48, 0.18);
  border-radius: 999px;
  cursor: pointer;
  font-weight: 700;
  font-size: 12px;
}

.today-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 4px 2px 12px;
  color: #3f4652;
  font-size: 12px;
  font-weight: 600;
}

.today-info p {
  margin: 0;
}

/* 左右布局 */
.main-layout {
  display: flex;
  gap: 20px;
}

.left-column {
  flex: 1;
  min-width: 0;
}

.right-column {
  width: 380px;
  flex-shrink: 0;
}

.section {
  background: rgba(255, 255, 255, 0.74);
  border: 1px solid rgba(255, 255, 255, 0.78);
  border-radius: 24px;
  box-shadow: 0 20px 55px rgba(31, 35, 40, 0.18), inset 0 1px 0 rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(18px) saturate(1.25);
  -webkit-backdrop-filter: blur(18px) saturate(1.25);
  padding: 16px;
  margin-bottom: 14px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: -2px -2px 14px;
  padding: 0 2px 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.72);
}

.section-header h2 {
  margin: 0;
  color: #1d1d1f;
  font-weight: 600;
  font-size: 16px;
}

.sub-title {
  color: #86868b;
  font-size: 13px;
}

.module-layer {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.module {
  position: relative;
  overflow: hidden;
  min-height: 98px;
  padding: 14px;
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.76);
  background: rgba(255, 255, 255, 0.66);
  color: #1d1d1f;
  text-align: left;
  cursor: pointer;
  box-shadow: 0 12px 32px rgba(31, 35, 40, 0.14);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
}

.module.active {
  background: rgba(0, 122, 255, 0.12);
  color: #007aff;
  border-color: rgba(0, 122, 255, 0.16);
  box-shadow: none;
}

.module h3 {
  margin: 0 0 8px;
  font-size: 17px;
}

.module p {
  margin: 0;
  color: inherit;
  opacity: 0.78;
  font-size: 12px;
  line-height: 1.45;
}

.module .num {
  position: absolute;
  right: 12px;
  bottom: 10px;
  font-size: 22px;
  font-weight: 800;
  opacity: 0.18;
}

/* 用户资料 */
.profile-section {
  padding: 16px;
}

.home-profile-section {
  max-width: none;
  margin-bottom: 18px;
}

.home-profile-section .profile-card {
  gap: 14px;
}

.home-profile-section .profile-top {
  min-height: 78px;
}

.home-profile-section .profile-points-cards {
  margin-top: 2px;
}

.profile-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.profile-top {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.profile-avatar {
  width: 56px;
  height: 56px;
  flex: 0 0 auto;
  border-radius: 18px;
  overflow: hidden;
  background: linear-gradient(135deg, #007aff, #34c759);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 800;
  box-shadow: 0 8px 22px rgba(0, 122, 255, 0.28);
}

.profile-avatar.large {
  width: 86px;
  height: 86px;
  font-size: 32px;
}

.profile-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-details {
  flex: 1;
  min-width: 0;
}

.profile-title {
  display: block;
  margin-bottom: 4px;
  color: #6e7781;
  font-size: 12px;
  font-weight: 600;
}

.profile-details h3 {
  margin: 0 0 4px;
  color: #1d1d1f;
  font-size: 19px;
  font-weight: 700;
}

.profile-details p {
  margin: 0;
  color: #6e7781;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.profile-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.profile-points-cards {
  width: 100%;
}

.edit-profile-btn {
  padding: 7px 10px;
  border: 1px solid rgba(0, 122, 255, 0.16);
  border-radius: 999px;
  background: rgba(0, 122, 255, 0.12);
  color: #007aff;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.profile-modal {
  max-width: 420px;
}

.avatar-edit {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 18px;
}

.avatar-upload-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 8px 12px;
  border: 1px solid rgba(0, 122, 255, 0.16);
  border-radius: 999px;
  background: rgba(0, 122, 255, 0.12);
  color: #007aff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
}

.avatar-upload-btn input {
  display: none;
}

.error {
  color: #ff3b30;
  font-size: 13px;
}

.success {
  color: #34c759;
  font-size: 13px;
}

/* 创建任务表单 */
.create-task form, .action-form form {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 12px;
  padding: 10px;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.48);
}

.create-task input, .create-task select,
.action-form input, .action-form select {
  padding: 9px 12px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 12px;
  font-size: 14px;
  color: #1d1d1f;
  background: rgba(255, 255, 255, 0.62);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.55);
}

.create-task input:focus, .create-task select:focus,
.action-form input:focus, .action-form select:focus {
  border-color: rgba(0, 122, 255, 0.55);
  outline: none;
  background: rgba(255, 255, 255, 0.88);
}

.create-task input:first-child { flex: 1; min-width: 140px; }
.action-form input[placeholder="描述"] { flex: 1; min-width: 100px; }

.list-controls {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.filter-input,
.sort-select {
  padding: 9px 12px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 12px;
  font-size: 14px;
  color: #1d1d1f;
  background: rgba(255, 255, 255, 0.62);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.55);
}

.filter-input:focus,
.sort-select:focus {
  border-color: rgba(0, 122, 255, 0.55);
  outline: none;
  background: rgba(255, 255, 255, 0.88);
}

.filter-input {
  flex: 1;
  min-width: 150px;
}

.sort-select {
  min-width: 120px;
}

.date-wheel {
  display: flex;
  gap: 6px;
  flex-wrap: nowrap;
  min-width: 210px;
}

.datetime-wheel {
  min-width: 330px;
}

.date-wheel select {
  min-width: 0;
  padding-right: 8px;
  appearance: auto;
}

.create-task button, .action-form button {
  padding: 9px 16px;
  background: rgba(0, 122, 255, 0.12);
  color: #007aff;
  border: 1px solid rgba(0, 122, 255, 0.16);
  border-radius: 999px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 700;
}

.create-task button:disabled, .action-form button:disabled {
  background: rgba(142, 142, 147, 0.12);
  color: #8e8e93;
  border-color: rgba(142, 142, 147, 0.16);
  cursor: not-allowed;
}

/* 任务列表 */
.tasks {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  margin-top: 12px;
  padding-top: 10px;
}

.page-btn {
  padding: 7px 16px;
  background: rgba(0, 122, 255, 0.12);
  color: #007aff;
  border: 1px solid rgba(0, 122, 255, 0.16);
  border-radius: 999px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 700;
}

.page-btn:disabled {
  background: rgba(142, 142, 147, 0.12);
  color: #8e8e93;
  border-color: rgba(142, 142, 147, 0.16);
  cursor: not-allowed;
}

.page-info {
  color: #86868b;
  font-size: 13px;
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  padding: 11px;
  background: rgba(255, 255, 255, 0.52);
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-left: 5px solid #34c759;
  border-radius: 16px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.task-item.inactive {
  opacity: 0.72;
}

.task-item.checked {
  background: rgba(255, 255, 255, 0.38);
}

.task-item.checked .task-info h4 {
  color: #86868b;
}

.task-item.checked .level-dot {
  background: #86868b;
}

.task-item.checked .task-meta {
  color: #86868b;
}

.task-item.checked .task-points {
  color: #86868b;
}

.task-item.level-1.checked { border-left-color: #86868b; }
.task-item.level-2.checked { border-left-color: #86868b; }
.task-item.level-3.checked { border-left-color: #86868b; }

.task-item.level-1 { border-left-color: #34c759; }
.task-item.level-2 { border-left-color: #ff9500; }
.task-item.level-3 { border-left-color: #ff3b30; }

.task-info h4 {
  margin: 0 0 4px 0;
  display: flex;
  align-items: center;
  gap: 6px;
  color: #1d1d1f;
  font-size: 14px;
}

.level-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.dot-1 { background: #34c759; }
.dot-2 { background: #ff9500; }
.dot-3 { background: #ff3b30; }

.task-meta {
  display: flex;
  gap: 8px;
  color: #86868b;
  font-size: 12px;
}

.task-points {
  color: #34c759;
  font-weight: 500;
}

.task-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

/* 手机端按钮文字竖排 */
@media (max-width: 600px) {
  .task-actions button,
  .task-item button {
    writing-mode: vertical-rl;
    padding: 6px 3px !important;
    font-size: 12px !important;
    letter-spacing: 0;
  }
}

.edit-btn,
.delete-btn,
.skip-btn,
.repeat-btn,
.task-item button:not(.delete-btn):not(.skip-btn):not(.edit-btn) {
  padding: 6px 12px;
  border-radius: 999px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 700;
}

.edit-btn,
.task-item button:not(.delete-btn):not(.skip-btn):not(.edit-btn):not(.cancel-btn) {
  background: rgba(0, 122, 255, 0.12);
  color: #007aff;
  border: 1px solid rgba(0, 122, 255, 0.16);
}

.delete-btn {
  background: rgba(255, 59, 48, 0.12);
  color: #d70015;
  border: 1px solid rgba(255, 59, 48, 0.18);
}

.skip-btn,
.task-item button.cancel-btn {
  background: rgba(142, 142, 147, 0.12);
  color: #6e6e73;
  border: 1px solid rgba(142, 142, 147, 0.18);
}

.repeat-btn {
  background: rgba(52, 199, 89, 0.12);
  color: #1f8f45;
  border: 1px solid rgba(52, 199, 89, 0.18);
}

.delete-btn.small { padding: 5px 9px; }

.task-item button:disabled {
  background: rgba(142, 142, 147, 0.12);
  color: #8e8e93;
  border-color: rgba(142, 142, 147, 0.16);
}

/* 积分历史 */
.history-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-item,
.long-term-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  padding: 11px;
  background: rgba(255, 255, 255, 0.52);
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 16px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.item-info {
  display: flex;
  gap: 8px;
  align-items: center;
}

.type {
  padding: 3px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.type.earn {
  background: #34c759;
  color: white;
}

.type.spend {
  background: #ff3b30;
  color: white;
}

.desc { color: #1d1d1f; font-size: 13px; }
.time { color: #86868b; font-size: 11px; }

.item-right {
  display: flex;
  gap: 8px;
  align-items: center;
}

.amount { font-weight: 600; font-size: 13px; }
.amount.earn { color: #34c759; }
.amount.spend { color: #ff3b30; }

.long-term-layout {
  width: 100%;
}

.long-term-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 12px;
}

.stat-card {
  padding: 12px;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.52);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.stat-card label {
  display: block;
  margin-bottom: 6px;
  color: #6e7781;
  font-size: 12px;
}

.stat-card strong {
  color: #1d1d1f;
  font-size: 18px;
}

.long-term-list-block {
  margin-top: 14px;
}

.inline-header {
  margin-top: 0;
}

.item-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item-icon {
  width: 38px;
  height: 38px;
  flex: 0 0 auto;
  border-radius: 12px;
  display: grid;
  place-items: center;
  background: #d6f5dd;
  color: #1d1d1f;
  font-weight: 800;
}

.long-term-item.scrapped {
  opacity: 0.74;
}

.long-term-main {
  flex: 1;
  min-width: 0;
}

.long-term-main h4 {
  margin: 0 0 4px;
  color: #1d1d1f;
  font-size: 14px;
}

.long-term-main p {
  margin: 0;
  color: #6e7781;
  font-size: 12px;
}

.long-term-side {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #3f4652;
  font-size: 12px;
  font-weight: 700;
}

.long-term-actions {
  display: flex;
  gap: 6px;
}

.long-term-edit-form {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  width: 100%;
}

.long-term-edit-form input,
.long-term-edit-form select {
  padding: 9px 12px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 12px;
  font-size: 14px;
  color: #1d1d1f;
  background: rgba(255, 255, 255, 0.62);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.55);
}

.long-term-edit-form input:first-child {
  flex: 1;
  min-width: 120px;
}

.long-term-edit-form input[type="number"] {
  width: 110px;
}

/* 图表 */
.days-select {
  padding: 6px 10px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 999px;
  font-size: 13px;
  color: #1d1d1f;
  background: rgba(255, 255, 255, 0.62);
}

.chart-wrapper {
  height: 300px;
  position: relative;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 18px;
  padding: 15px;
  background: rgba(255, 255, 255, 0.46);
}

.chart-wrapper canvas {
  width: 100% !important;
  height: 100% !important;
}

.heatmap-card {
  margin-bottom: 18px;
}

.heatmap-header {
  gap: 12px;
}

.heatmap-scroll {
  overflow: visible;
}

.heatmap-content {
  width: 100%;
  min-width: 100%;
}

.heatmap-months {
  display: grid;
  gap: 4px;
  margin-bottom: 5px;
}

.heatmap-month-label {
  color: #57606a;
  font-size: 11px;
  line-height: 12px;
  white-space: nowrap;
}

.heatmap-board {
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 18px;
  padding: 14px;
  background: rgba(255, 255, 255, 0.46);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.5);
}

.heatmap-grid {
  display: grid;
  gap: clamp(1px, 0.7vw, 4px);
  width: 100%;
}

.heatmap-weekdays,
.heatmap-week-column {
  display: grid;
  grid-template-rows: repeat(7, clamp(9px, 2.4vw, 13px));
  gap: clamp(1px, 0.7vw, 4px);
}

.heatmap-weekdays span {
  color: #57606a;
  font-size: 10px;
  line-height: 13px;
}

.heatmap-cell {
  display: block;
  width: 11px;
  height: 11px;
  flex: 0 0 auto;
  border-radius: 3px;
  background: #ebedf0;
}

.heatmap-day-cell {
  position: relative;
  width: 100%;
  min-width: 0;
  height: clamp(9px, 2.4vw, 13px);
  border: 1px solid rgba(27, 31, 36, 0.06);
  cursor: pointer;
}

.heatmap-tooltip {
  position: absolute;
  left: 50%;
  bottom: calc(100% + 8px);
  z-index: 20;
  transform: translate3d(-50%, 0, 0);
  min-width: 120px;
  padding: 7px 9px;
  border-radius: 8px;
  background: rgba(29, 29, 31, 0.92);
  color: #fff;
  font-size: 11px;
  line-height: 1.5;
  text-align: left;
  white-space: pre-line;
  pointer-events: none;
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.18);
  will-change: transform, opacity;
}

.heatmap-tooltip::after {
  content: '';
  position: absolute;
  left: 50%;
  top: 100%;
  transform: translateX(-50%);
  border: 6px solid transparent;
  border-top-color: rgba(29, 29, 31, 0.92);
}

.heatmap-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 4px;
  margin-top: 10px;
  color: #57606a;
  font-size: 11px;
}

.heatmap-cell.zero,
.heatmap-cell.heatmap-empty {
  background: #ebedf0;
}

.heatmap-cell.pos-1 { background: #d6f5dd; }
.heatmap-cell.pos-2 { background: #9be9a8; }
.heatmap-cell.pos-3 { background: #6ccf7f; }
.heatmap-cell.neg-1 { background: #ffd9d6; }
.heatmap-cell.neg-2 { background: #ffb4ad; }
.heatmap-cell.neg-3 { background: #ff8a83; }

@media (max-width: 600px) {
  .dashboard-container {
    width: 100%;
    max-width: 100vw;
    overflow-x: hidden;
    padding: 15px 8px 50px 8px;
  }

  .profile-card {
    gap: 10px;
  }

  .profile-top {
    align-items: flex-start;
  }

  .profile-actions {
    flex-direction: column;
    gap: 6px;
  }

  .profile-avatar {
    width: 46px;
    height: 46px;
    font-size: 18px;
  }

  .edit-profile-btn {
    padding: 6px 8px;
    font-size: 11px;
  }

  .avatar-edit {
    gap: 12px;
  }

  .profile-points-cards {
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 5px;
  }

  .point-card {
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 2px;
    min-width: 0;
    min-height: 46px;
    padding: 7px 3px;
    font-size: 10px;
    line-height: 1.15;
    white-space: normal;
    overflow-wrap: anywhere;
  }

  .point-label,
  .point-value {
    display: block;
    min-width: 0;
  }

  .point-label {
    font-size: 10px;
  }

  .point-value {
    font-size: 11px;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .point-card.daily-cost {
    grid-column: 1 / -1;
    min-height: 38px;
    flex-direction: row;
    align-items: center;
  }

  .heatmap-header {
    flex-direction: column;
    gap: 8px;
  }

  .section {
    padding: 12px;
    border-radius: 18px;
  }

  .module-layer,
  .long-term-summary {
    grid-template-columns: 1fr;
  }

  .date-wheel,
  .datetime-wheel {
    min-width: 100%;
    flex-wrap: wrap;
  }

  .date-wheel select {
    flex: 1;
  }

  .long-term-item,
  .history-item {
    align-items: flex-start;
  }

  .long-term-side {
    flex-direction: column;
    align-items: flex-end;
  }

  .heatmap-board {
    padding: 10px;
  }

  .heatmap-grid {
    gap: 2px;
  }

  .heatmap-weekdays,
  .heatmap-week-column {
    grid-template-rows: repeat(7, clamp(7px, 2.2vw, 10px));
    gap: 2px;
  }

  .heatmap-day-cell {
    height: clamp(7px, 2.2vw, 10px);
  }

  .heatmap-weekdays span {
    font-size: 9px;
    line-height: clamp(7px, 2.2vw, 10px);
  }
}

.loading, .empty {
  text-align: center;
  padding: 20px;
  color: #86868b;
}

/* 响应式 */
@media (max-width: 900px) {
  .main-layout {
    flex-direction: column;
  }
  .right-column {
    width: 100%;
  }
}

/* 编辑弹窗 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.38);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(6px);
  -webkit-backdrop-filter: blur(6px);
}

.modal-content {
  background: rgba(255, 255, 255, 0.84);
  border: 1px solid rgba(255, 255, 255, 0.78);
  padding: 22px;
  border-radius: 24px;
  width: 350px;
  box-shadow: 0 20px 55px rgba(31, 35, 40, 0.22), inset 0 1px 0 rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(18px) saturate(1.25);
  -webkit-backdrop-filter: blur(18px) saturate(1.25);
}

.modal-content h3 {
  margin: 0 0 20px 0;
  color: #1d1d1f;
  font-size: 18px;
}

.modal-content .form-group {
  margin-bottom: 15px;
}

.modal-content label {
  display: block;
  margin-bottom: 5px;
  color: #1d1d1f;
  font-size: 14px;
}

.modal-content input,
.modal-content select {
  width: 100%;
  padding: 10px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 12px;
  font-size: 14px;
  color: #1d1d1f;
  background: rgba(255, 255, 255, 0.62);
  box-sizing: border-box;
}

.modal-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.modal-actions button {
  flex: 1;
  padding: 12px;
  border-radius: 999px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 700;
}

.modal-actions .cancel-btn {
  background: rgba(142, 142, 147, 0.12);
  color: #6e6e73;
  border: 1px solid rgba(142, 142, 147, 0.18);
}

.modal-actions .save-btn {
  background: rgba(0, 122, 255, 0.12);
  color: #007aff;
  border: 1px solid rgba(0, 122, 255, 0.16);
}

/* 底部版权 */
.footer {
  text-align: center;
  padding: 15px 20px;
  color: #86868b;
  font-size: 12px;
  display: flex;
  justify-content: center;
  gap: 20px;
}

.footer a {
  color: #007aff;
  text-decoration: none;
}
</style>