#! /bin/bash

# Invoke Policy 01 (FAIL Scenario) - Ensure atleast 10 seconds is set for startingDeadlineSeconds while creating a CronJob
cd 01-startingDeadlineSeconds && datree publish policy-startingDeadlineSeconds.yaml && datree test fail-startingDeadlineSeconds.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 02 (FAIL Scenario) - Check if unmanaged jobs are not left around after it's fully deleted
cd .. && cd 02-ttlSecondsAfterFinished && datree publish policy-ttlSecondsAfterFinished.yaml && datree test fail-ttlSecondsAfterFinished.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 03 (FAIL Scenario) - Provide atleast 1 revisionHistoryLimit to ensure successful rollback of deployment
cd .. && cd 03-revisionHistoryLimit && datree publish policy-revisionHistoryLimit.yaml && datree test fail-revisionHistoryLimit.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 04 (FAIL Scenario) - Ensure a proper namespace exists to prevent accidentlal creation of pods in the active namespace
cd .. && cd 04-namespacePolicy && datree publish policy-namespacePolicy.yaml && datree test fail-namespacePolicy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 05 (FAIL Scenario) - minAvailable value set for PDB to ensure that the number of replicas running is never brought below the number needed
cd .. && cd 05-podDisruptionBudget && datree publish policy-podDisruptionBudget.yaml && datree test fail-podDisruptionBudget.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 06 (FAIL Scenario) - Add a properly configured startupProbe to ensure app within the container is started properly
cd .. && cd 06-startupProbe && datree publish policy-startupProbe.yaml && datree test fail-startupProbe.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 07 (FAIL Scenario) - Avoid using the :latest tag for images in prod as its hard to track image version used
cd .. && cd 07-imageAvoidLatest && datree publish policy-imageAvoidLatest.yaml && datree test fail-imageAvoidLatest.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 08 (FAIL Scenario) - Use only the accepted ipFamilyPolicy options
cd .. && cd 08-ipFamilyPolicy && datree publish policy-ipFamilyPolicy.yaml && datree test fail-ipFamilyPolicy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 09 (FAIL Scenario) - Use only the accepted concurrencyPolicy options
cd .. && cd 09-concurrencyPolicy && datree publish policy-concurrencyPolicy.yaml && datree test fail-concurrencyPolicy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 10 (FAIL Scenario) - Use only the accepted LimitRange type options
cd .. && cd 10-limitRangeType && datree publish policy-limitRangeType.yaml && datree test fail-limitRangeType.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 11 (FAIL Scenario) - Use only the accepted strategy options
cd .. && cd 11-deploymentStrategy && datree publish policy-deploymentStrategy.yaml && datree test fail-deploymentStrategy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 12 (FAIL Scenario) - Use only the accepted imagePullPolicy options
cd .. && cd 12-imagePullPolicy && datree publish policy-imagePullPolicy.yaml && datree test fail-imagePullPolicy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 13 (FAIL Scenario) - Add proper secrets to make sure the image pull happens without failure while using private registry
cd .. && cd 13-imagePullSecrets && datree publish policy-imagePullSecrets.yaml && datree test fail-imagePullSecrets.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 14 (FAIL Scenario) - Add proper SELinux level for your container
cd .. && cd 14-seLinuxOptions && datree publish policy-seLinuxOptions.yaml && datree test fail-seLinuxOptions.yaml --ignore-missing-schemas && cd ..
