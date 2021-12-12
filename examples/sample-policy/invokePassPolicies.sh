#! /bin/bash

# Invoke Policy 01 (PASS Scenario) - Ensure atleast 10 seconds is set for startingDeadlineSeconds while creating a CronJob
cd 01-startingDeadlineSeconds && datree publish policy-startingDeadlineSeconds.yaml && datree test pass-startingDeadlineSeconds.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 02 (PASS Scenario) - Check if unmanaged jobs are not left around after its fully deleted
cd 02-ttlSecondsAfterFinished\ && datree publish policy-ttlSecondsAfterFinished.yaml && datree test pass-ttlSecondsAfterFinished.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 03 (PASS Scenario) - Provide atleast 1 revisionHistoryLimit to ensure successful rollback of deployment
cd 03-revisionHistoryLimit && datree publish policy-revisionHistoryLimit.yaml && datree test pass-revisionHistoryLimit.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 04 (PASS Scenario) - Ensure a proper namespace exists to prevent accidentlal creation of pods in the active namespace
cd 04-namespacePolicy && datree publish policy-namespacePolicy.yaml && datree test pass-namespacePolicy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 05 (PASS Scenario) - minAvailable value set for PDB to ensure that the number of replicas running is never brought below the number needed
cd 05-podDisruptionBudget && datree publish policy-podDisruptionBudget.yaml && datree test pass-podDisruptionBudget.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 06 (PASS Scenario) - Add a properly configured startupProbe to ensure app within the container is started properly
cd 06-startupProbe && datree publish policy-startupProbe.yaml && datree test pass-startupProbe.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 07 (PASS Scenario) - Avoid using the :latest tag for images in prod as its hard to track image version used
cd 07-imageAvoidLatest && datree publish policy-imageAvoidLatest.yaml && datree test pass-imageAvoidLatest.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 08 (PASS Scenario) - Use only the accepted ipFamilyPolicy options
cd 08-ipFamilyPolicy && datree publish policy-ipFamilyPolicy.yaml && datree test pass-ipFamilyPolicy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 09 (PASS Scenario) - Use only the accepted concurrencyPolicy options
cd 09-concurrencyPolicy && datree publish policy-concurrencyPolicy.yaml && datree test pass-concurrencyPolicy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 10 (PASS Scenario) - Use only the accepted LimitRange type options
cd 10-limitRangeType && datree publish policy-limitRangeType.yaml && datree test pass-limitRangeType.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 11 (PASS Scenario) - Use only the accepted strategy options
cd 11-deploymentStrategy && datree publish policy-deploymentStrategy.yaml && datree test pass-deploymentStrategy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 12 (PASS Scenario) - Use only the accepted imagePullPolicy options
cd 12-imagePullPolicy && datree publish policy-imagePullPolicy.yaml && datree test pass-imagePullPolicy.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 13 (PASS Scenario) - Add proper secrets to make sure the image pull happens without failure while using private registry
cd 13-imagePullSecrets && datree publish policy-imagePullSecrets.yaml && datree test pass-imagePullSecrets.yaml --ignore-missing-schemas && cd ..

# Invoke Policy 14 (PASS Scenario) - Add proper SELinux level for your container
cd 14-seLinuxOptions && datree publish policy-seLinuxOptions.yaml && datree test pass-seLinuxOptions.yaml --ignore-missing-schemas && cd ..
