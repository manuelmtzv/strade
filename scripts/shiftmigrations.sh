#!/bin/bash

POS=$1
SHIFT=$2
MIGRATIONS_PATH="./cmd/migrate/migrations"

if [[ -z "$POS" || -z "$SHIFT" ]]; then
  echo "Uso: ./shift_migrations.sh <posicion> <cantidad>"
  exit 1
fi


if ! [[ "$POS" =~ ^[0-9]+$ && "$SHIFT" =~ ^[0-9]+$ ]]; then
  echo "Error: La posición y el desplazamiento deben ser números enteros."
  exit 1
fi

find "$MIGRATIONS_PATH" -type f -printf "%f\n" | \
  grep -E '^[0-9]{6}_.+\.(up|down)\.sql$' | \
  sort -r | while read FILE; do
    NUM=$(echo $FILE | cut -d_ -f1)
    REST=$(echo $FILE | cut -d_ -f2-)

    NUM_INT=$(echo $NUM | sed 's/^0*//')
    if [[ "$NUM_INT" -gt "$POS" ]]; then
      NEW_NUM=$(printf "%06d" $((NUM_INT + SHIFT)))
      mv "$MIGRATIONS_PATH/$FILE" "$MIGRATIONS_PATH/${NEW_NUM}_$REST"
      echo "Renombrado: $FILE -> ${NEW_NUM}_$REST"
    fi
done