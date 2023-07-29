## Анализ требований

### Авторизация в системе
- Actor - User
- Command - Login
- Data - PopugBeak
- Event - User.Authorized

### Создание задачи  
- Actor - User
- Command - CreateTask
- Data - Task + CreatorUser.Id + AssignedUser.ID
- Event - Task.Assigned

### Изменение статуса задачи
- Actor - User
- Command - ChangeTaskStatus
- Data - Task.ID + NewStatus + ExecutorUser.ID
- Event - Task.StatusChanged

### Заасайнить все задачи
- Actor - Admin/Manager User (user with role admin or manager)
- Command - AssigneeTasks
- Data - ?
- Event - Task.Assigned

### Списание средств за назначенную задачу
- Actor - event Task.Assigned
- Command - ChargeMoneyForTask
- Data - User.ID + Task.ID
- Event - Balance.Changed

### Начисление средств за закрытую задачу
- Actor - event Task.StatusChanged.ClosedStatus
- Command - AccrueMoneyForTask
- Data - User.ID + Task.ID
- Event - Balance.Changed

### Закрытие дня
- Actor - event.EndOfTheDay
- Command - EndOfTheDay
- Data - ? maybe current date
- Event - Day.Closed

### Подсчёт выручки пользователя за день
- Actor - event Day.Closed
- Command - CalcEarnedMoney
- Data - User.ID
- Event - Money.Paid

### Отслеживание изменений баланса
- Actor - event Money.Paid / Balance.Changed
- Command - CreateAccountingAuditLogRecord
- Data - Balance.ChangeAmount + BalanceChangeInfo(salary/task)
- Event - AuditLogRecord.Created

### Получение аналитики
- Actor - User.Admin
- Command - GetStats
- Data - CurrentDate
- Event - ?
