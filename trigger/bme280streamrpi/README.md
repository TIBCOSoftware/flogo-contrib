# bme280streamrpi


Trigger for sensor BME280 for Raspberry Pi (Temperature, Humidity, Pressure)

## Installation

#### Install Trigger
Example: install **bme280streamrpi** trigger

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/trigger/bme280streamrpi
```


## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "delay_ms",
      "type": "integer",
      "required": false
    }
  ],
  "outputs": [
   {
      "name": "Temperature",
      "type": "number"
    },
    {
      "name": "Pressure",
      "type": "number"
    },
    {
      "name": "Humidity",
      "type": "number"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| delay_ms      | The delay, in milliseconds, between two measures |         
Note: |* **delay_ms**: If left blank, defaut value of 500ms



## Output
| Setting     | Description    |
|:------------|:---------------|
| Temperature      | The temperature, in degree celsius |         
| Pressure      | The pressure, in hPa |      
| Humidity      | The humidity, in percentage |     

