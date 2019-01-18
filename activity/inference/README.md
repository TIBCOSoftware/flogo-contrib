---
title: ML Model Inference
weight: 4618
---

# Inference
This activity enables the inferencing of Machine Learning models within Flogo applications. This activity was built with a framework contribution model concept. The implemented framework is TensorFlow.

For detailed instructions, refer to the [Flogo Documentation](https://tibcosoftware.github.io/flogo/development/flows/tensorflow/).

## Installation
### Flogo Web
This activity does not come pre-installed with the Web UI for a number of reasons, such as, the size, the requirement of TensorFlow lib and also the fact that it is an activity that is commonly used on a daily basis.

### Flogo CLI
```bash
flogo install activity github.com/TIBCOSoftware/flogo-contrib/activity/inference
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "model",
      "type": "string",
      "required": true
    },
    {
      "name": "framework",
      "type": "string",
      "required": true
    },
    {
      "name": "sigDefName",
      "type": "string",
      "required": false,
      "value":"serving_default"
    },
    {
      "name": "tag",
      "type": "string",
      "required": false,
      "value":"serve"
    },
    {
      "name": "features",
      "type": "array",
      "required": true
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "object"
    }
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| model      | True     | The location to the archive. If using TensorFlow, the archive must contain the TensorFlow SavedModel |
| framework         | True     | The framework to use. Other frameworks can be registered at build time, the only available framework is `TensorFlow` |
| sigDefName       | False    | The default signature definition. This comes from the SavedModel metadata. The default value is `serving_default` |
| tag  | False    | The model tag. This comes from the SavedModel metadata. The default value is `serve` |
| features | true    | An array of input features. Refer to the following sample. |


## Example

### Estimators
The following example demonstrates how to invoke the inference activity and pass the input feature set for the tensor named `inputs`. Tensor names may vary, it is best to refer to the SavedModel metadata to identify the correct tensor name.

```json
{
  "id": "inference_2",
  "name": "Invoke ML Model",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/inference",
    "input": {
      "model": "Archive.zip",
      "framework": "Tensorflow"
    },
    "mappings": {
      "input": [
        {
          "type": "assign",
          "value": [
            {
              "name": "inputs",
              "data": {
                "z-axis-q75": 4.140586,
                "corr-x-z": 0.1381063882214782,
                "x-axis-mean": 1.7554575428900194,
                "z-axis-sd": 4.6888631696380765,
                "z-axis-skew": -0.3619011587545954,
                "y-axis-sd": -7.959084724314854,
                "y-axis-q75": 16.467001,
                "corr-z-y": 0.3467060369518231,
                "x-axis-sd": 6.450293741961166,
                "x-axis-skew": 0.09756801680727022,
                "y-axis-mean": 9.389463650669393,
                "y-axis-skew": -0.49036224958471764,
                "z-axis-mean": 1.1226106985139188,
                "x-axis-q25": -3.1463003,
                "x-axis-q75": 6.3198414,
                "y-axis-q25": 3.0645783,
                "z-axis-q25": -1.9477097,
                "corr-x-y": 0.08100326860866637
              }
            }
          ],
          "mapTo": "features"
        }
      ]
    }
  }
}
```