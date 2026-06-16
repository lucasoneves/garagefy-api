# Changelog - 15/06/2026

## Controllers

### vehicle.go
- Adicionado `Color` e `CurrentOdo` ao input struct do `CreateVehicle`
- Campos agora são passados para a criação do veículo

### service.go
- Adicionado `ShopName` e `CurrentOdo` ao input struct do `CreateService`
- Campos agora são passados para a criação do serviço

### logbook.go
- Adicionado `Category` ao input struct do `CreateLogbookEntry` (era `binding:"required"`)
- Adicionado `Category` ao input struct do `UpdateLogbookEntry`
- Ambos convertem a string para `models.LogbookCategory`

## Models

### logbook.go
- Removido campo `AttachmentURL *string` e sua tag

## Main

### main.go
- Removido `r.Static("/uploads", "./uploads")` (não mais necessário)
- Adicionado `gin.SetMode(gin.ReleaseMode)`

## Frontend (garagefy)

### logbook/page.tsx
- Removido `attachment_url` da interface `LogbookEntry`

### next.config.ts
- Removido `remotePattern` do `localhost:8080/uploads`
