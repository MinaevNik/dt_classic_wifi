#!/bin/bash

# Устанавливаем имя интерфейса вручную, так как оно уже известно
wifi_interface="wlx00367608058b"

# Проверка наличия интерфейса
if [ -z "$wifi_interface" ]; then
  echo "Не найден интерфейс беспроводной сети. Проверьте, что ваш беспроводной адаптер правильно установлен."
  exit 1
fi

echo "Обнаружен интерфейс: $wifi_interface"

# Включение интерфейса
sudo ip link set $wifi_interface up

# Сканируем доступные сети и выводим их в список
echo "Сканирование доступных сетей..."
sudo nmcli dev wifi rescan
networks=$(sudo nmcli -t -f SSID dev wifi list)

if [ -z "$networks" ]; then
  echo "Не удалось найти доступные сети. Проверьте, что ваш беспроводной адаптер включен и находится в зоне действия сети."
  exit 1
fi

# Выводим список сетей для выбора пользователем
echo "Доступные сети:"
IFS=$'\n'
networks_array=($networks)
for i in "${!networks_array[@]}"; do
  echo "$i) ${networks_array[$i]}"
done

# Просим пользователя выбрать сеть
echo "Выберите номер сети:"
read network_number

# Проверяем правильность ввода
if ! [[ "$network_number" =~ ^[0-9]+$ ]] || [ "$network_number" -ge "${#networks_array[@]}" ]]; then
  echo "Некорректный ввод. Пожалуйста, запустите скрипт снова и выберите правильный номер сети."
  exit 1
fi

selected_ssid=${networks_array[$network_number]}

# Просим пользователя ввести пароль
echo "Введите пароль для сети '$selected_ssid':"
read -s network_password

# Подключаемся к выбранной сети
sudo nmcli dev wifi connect "$selected_ssid" password "$network_password" ifname $wifi_interface

echo "Подключение к сети '$selected_ssid' выполнено."
